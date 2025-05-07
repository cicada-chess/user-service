// cmd/user-service/main.go
package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	pb "gitlab.mai.ru/cicada-chess/backend/auth-service/pkg/auth"
	profileService "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/profile"
	userService "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/config"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/db/minio"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/db/postgres"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/messaging/kafka"
	profileStorage "gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/minio/profile"
	profileInfrastructure "gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/postgres/profile"
	userInfrastructure "gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/postgres/user"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/grpc/handlers"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp"
	"gitlab.mai.ru/cicada-chess/backend/user-service/logger"
	"gitlab.mai.ru/cicada-chess/backend/user-service/pkg/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// @title User API
// @version 1.0
// @description API для управления пользователями

// @host cicada-chess.ru:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	log := logger.New()

	conn, err := grpc.NewClient("auth-service:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	config, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	client := pb.NewAuthServiceClient(conn)

	dbConn, err := postgres.NewPostgresDB(config.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	storageConn, err := minio.NewMinioStorage(config.Storage)
	if err != nil {
		log.Fatalf("Failed to connect to storage: %v", err)
	}

	userRepo := userInfrastructure.NewUserRepository(dbConn)
	profileRepo := profileInfrastructure.NewProfileRepository(dbConn)
	profileStorage := profileStorage.NewProfileStorage(storageConn, config.Storage.BucketName, config.Storage.Host)

	kafkaProducer, err := kafka.NewKafkaProducer(config.Kafka.Brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	notificationSender := kafka.NewKafkaNotificationSender(kafkaProducer, config.Kafka.Topic, log)

	userService := userService.NewUserService(userRepo, notificationSender)

	profileService := profileService.NewProfileService(profileRepo, userRepo, profileStorage, client)

	r := gin.Default()
	ginapp.InitRoutes(r, userService, profileService, log)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	grpcServer := grpc.NewServer()
	grpcHandler := handlers.NewGRPCHandler(userService, profileService)
	user.RegisterUserServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("Failed to listen on :9090: %v", err)
		}
		log.Println("Starting gRPC server on :9090")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	grpcServer.GracefulStop()

	log.Println("Server stopped")
}

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
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/db/postgres"
	profileInfrastructure "gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/postgres/profile"
	userInfrastructure "gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/postgres/user"
	profileStorage "gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/storage/profile"
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

// @host 217.114.11.158:8080
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

	client := pb.NewAuthServiceClient(conn)

	cfgToDB := postgres.GetDBConfig()
	dbConn, err := postgres.NewPostgresDB(cfgToDB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	userRepo := userInfrastructure.NewUserRepository(dbConn)
	profileRepo := profileInfrastructure.NewProfileRepository(dbConn)
	profileStorage := profileStorage.NewProfileStorage("/uploads/avatars")

	userService := userService.NewUserService(userRepo)

	profileService := profileService.NewProfileService(profileRepo, userRepo, profileStorage, client)

	r := gin.Default()
	ginapp.InitRoutes(r, userService, profileService, log)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	grpcServer := grpc.NewServer()
	grpcHandler := handlers.NewGRPCHandler(userService)
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

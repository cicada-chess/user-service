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
	service "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/db/postgres"
	infrastructure "gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/postgres/user"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/grpc/handlers"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp"
	"gitlab.mai.ru/cicada-chess/backend/user-service/logger"
	"gitlab.mai.ru/cicada-chess/backend/user-service/pkg/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// @title User API
// @version 1.0
// @description API для управления пользователями

// @host localhost:8080
// @BasePath /

func main() {
	log := logger.New()

	cfgToDB := postgres.GetDBConfig()
	dbConn, err := postgres.NewPostgresDB(cfgToDB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	userRepo := infrastructure.NewUserRepository(dbConn)

	userService := service.NewUserService(userRepo)

	r := gin.Default()
	ginapp.InitRoutes(r, userService, log)

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

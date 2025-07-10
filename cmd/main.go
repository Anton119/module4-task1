package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"module4-task1/internal/logger"
	"net"
	"os"

	"module4-task1/internal/api"
	auth2 "module4-task1/internal/api/proto/auth"
	"module4-task1/internal/repo"
	"module4-task1/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func buildDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	dsn := buildDSN()

	db, err := repo.NewPostgres(dsn)
	if err != nil {
		logger.Log.Info("Failed to connect to database", zap.Error(err))
	}

	userRepo := repo.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		log.Fatal("GRPC_PORT environment variable not set")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	authServer := &api.GrpcAuthServer{Service: authService}
	auth2.RegisterAuthServiceServer(grpcServer, authServer)

	reflection.Register(grpcServer)

	logger.Log.Info("Starting server")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Log.Fatal("Error starting server", zap.Error(err))
	}
}

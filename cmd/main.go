package main

import (
	"log"
	"net"
	"os"

	"module4-task1/internal/api"
	auth2 "module4-task1/internal/api/proto/auth"
	"module4-task1/internal/repo"
	"module4-task1/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Получаем DSN из переменной окружения или укажи напрямую
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://anton:secret@localhost:5432/users?sslmode=disable"
	}

	// Подключаемся к БД
	db, err := repo.NewPostgres(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Создаём репозиторий пользователей
	userRepo := repo.NewUserRepository(db)

	// Создаём сервис аутентификации
	authService := service.NewAuthService(userRepo)

	// Настраиваем gRPC сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Регистрируем gRPC сервер с реализацией, которая использует сервис
	authServer := &api.GrpcAuthServer{Service: authService}
	auth2.RegisterAuthServiceServer(grpcServer, authServer)

	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

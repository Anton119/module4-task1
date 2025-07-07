package main

import (
	"context"
	"log"
	auth2 "module4-task1/internal/api/proto/auth"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Реализация сервера
type authServer struct {
	auth2.UnimplementedAuthServiceServer
}

// Реализация метода Register
func (s *authServer) Register(ctx context.Context, req *auth2.RegisterRequest) (*auth2.RegisterResponse, error) {
	log.Printf("Register called with username: %s", req.GetUsername())
	return &auth2.RegisterResponse{
		Message: "User " + req.GetUsername() + " registered successfully",
	}, nil
}

// Реализация метода Login
func (s *authServer) Login(ctx context.Context, req *auth2.LoginRequest) (*auth2.LoginResponse, error) {
	log.Printf("Login called with username: %s", req.GetUsername())
	return &auth2.LoginResponse{
		Token: "fake-jwt-token-for-" + req.GetUsername(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	auth2.RegisterAuthServiceServer(grpcServer, &authServer{})

	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

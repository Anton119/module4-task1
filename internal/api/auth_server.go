package api

import (
	"context"
	"go.uber.org/zap"
	"module4-task1/internal/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"module4-task1/internal/api/proto/auth"
	"module4-task1/internal/service"
)

type GrpcAuthServer struct {
	auth.UnimplementedAuthServiceServer
	Service *service.AuthService
}

func (s *GrpcAuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	err := s.Service.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {

		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &auth.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *GrpcAuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	accessToken, refreshToken, err := s.Service.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		logger.Log.Info("Failed to login", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &auth.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *GrpcAuthServer) Validate(ctx context.Context, req *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	return s.Service.Validate(ctx, req)
}

func (s *GrpcAuthServer) Refresh(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	return s.Service.Refresh(ctx, req)
}

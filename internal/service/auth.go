package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"module4-task1/internal/api/proto/auth"
	"module4-task1/internal/repo"
)

type AuthService struct {
	userRepo *repo.UserRepository
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(userRepo *repo.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Регистрация — чистая бизнес-логика с простыми аргументами
func (s *AuthService) Register(ctx context.Context, username, password string) error {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, repo.ErrUserNotFound) {
		return err
	}
	if user.ID != "" {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.CreateUser(ctx, repo.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		Role:         "user",
	})
}

// Логин — бизнес-логика возвращает JWT
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user.ID == "" {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Внутренняя бизнес-логика для проверки токена и возврата userID
func (s *AuthService) ValidateToken(token string) (string, error) {
	claims, err := ValidateJWT(token)
	if err != nil {
		return "", err
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token payload")
	}
	return sub, nil
}

// gRPC метод Validate — принимает protobuf ValidateRequest, возвращает ValidateResponse
func (s *AuthService) Validate(ctx context.Context, req *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	userID, err := s.ValidateToken(req.GetToken())
	if err != nil {
		return nil, err
	}
	return &auth.ValidateResponse{UserId: userID}, nil
}

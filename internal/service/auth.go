package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"module4-task1/internal/api/proto/auth"
	"module4-task1/internal/logger"
	"module4-task1/internal/repo"
	"time"
)

type AuthService struct {
	userRepo *repo.UserRepository
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(userRepo *repo.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, repo.ErrUserNotFound) {
		logger.Log.Info("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return err
	}
	if user.ID != "" {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Info("Failed to hash password", zap.String("password", password), zap.Error(err))
		return err
	}

	return s.userRepo.CreateUser(ctx, repo.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		Role:         "user",
	})
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Log.Info("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return "", "", err
	}
	if user.ID == "" {
		logger.Log.Info("Failed to get user by username", zap.String("username", username))
		return "", "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		logger.Log.Info("Failed to compare password", zap.String("password", password), zap.Error(err))
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := GenerateJWT(user.ID, user.Username, user.Role, 15*time.Minute)
	if err != nil {
		logger.Log.Info("Failed to generate access token", zap.String("username", username), zap.Error(err))
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		logger.Log.Info("Failed to generate refresh token", zap.String("username", username), zap.Error(err))
		return "", "", err
	}

	err = s.userRepo.UpdateRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		logger.Log.Info("Failed to update refresh token", zap.String("username", username), zap.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ValidateToken(token string) (string, error) {
	claims, err := ValidateJWT(token)
	if err != nil {
		logger.Log.Info("Failed to validate token", zap.String("token", token), zap.Error(err))
		return "", err
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		logger.Log.Info("Failed to validate token", zap.String("token", token))
		return "", errors.New("invalid token payload")
	}
	return sub, nil
}

func (s *AuthService) Validate(ctx context.Context, req *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	userID, err := s.ValidateToken(req.GetToken())
	if err != nil {
		logger.Log.Info("Failed to validate token", zap.String("token", req.GetToken()), zap.Error(err))
		return nil, err
	}
	return &auth.ValidateResponse{UserId: userID}, nil
}

func (s *AuthService) Refresh(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	user, err := s.userRepo.GetUserByRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, err := GenerateJWT(user.ID, user.Username, user.Role, 15*time.Minute)
	if err != nil {
		logger.Log.Info("Failed to generate access token", zap.String("username", user.Username), zap.Error(err))
		return nil, err
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		logger.Log.Info("Failed to generate refresh token", zap.String("username", user.Username), zap.Error(err))
		return nil, err
	}

	err = s.userRepo.UpdateRefreshToken(ctx, user.ID, newRefreshToken)
	if err != nil {
		logger.Log.Info("Failed to update refresh token", zap.String("username", user.Username), zap.Error(err))
		return nil, err
	}

	return &auth.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Вспомогательная функция генерации случайного refresh токена
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logger.Log.Info("Failed to generate refresh token", zap.Error(err))
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

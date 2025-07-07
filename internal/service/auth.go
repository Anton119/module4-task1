package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"module4-task1/internal/repo"
)

type AuthService struct {
	userRepo *repo.UserRepository
}

func NewAuthService(userRepo *repo.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

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

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user.ID == "" {
		return "", errors.New("user not found")
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return "token123", nil
}

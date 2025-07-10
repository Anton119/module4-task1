package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"module4-task1/internal/logger"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	FirstName    sql.NullString
	LastName     sql.NullString
	IsActive     bool
	Role         string
	LastLoginAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RefreshToken string
}

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user User) error {
	query := `
		INSERT INTO users 
			(email, username, password_hash, first_name, last_name, is_active, role)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Conn.ExecContext(ctx, query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.Role,
	)
	if err != nil {
		logger.Log.Error("Error inserting user", zap.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}
	logger.Log.Info("User created", zap.String("email", user.Email))
	return nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var user User
	query := `
		SELECT 
			id, email, username, password_hash, first_name, last_name, is_active, role, last_login_at, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	err := r.db.Conn.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.Role,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		logger.Log.Error("Error getting user by username", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info("User not found", zap.String("username", username))
			return User{}, ErrUserNotFound
		}
		logger.Log.Error("Error getting user by username", zap.Error(err))
		return User{}, fmt.Errorf("failed to get user by username: %w", err)
	}
	logger.Log.Info("Successfully got user by username", zap.String("username", username))
	return user, nil
}

func (r *UserRepository) UpdateRefreshToken(ctx context.Context, userID, refreshToken string) error {
	query := `
		UPDATE users
		SET refresh_token = $1, updated_at = now()
		WHERE id = $2
	`
	_, err := r.db.Conn.ExecContext(ctx, query, refreshToken, userID)
	if err != nil {
		logger.Log.Error("Error updating refresh token", zap.Error(err))
		return fmt.Errorf("failed to update refresh token: %w", err)
	}
	logger.Log.Info("Successfully updated refresh token", zap.String("refresh_token", refreshToken))
	return nil
}

func (r *UserRepository) GetUserByRefreshToken(ctx context.Context, refreshToken string) (User, error) {
	var user User
	query := `
		SELECT 
			id, email, username, password_hash, first_name, last_name, is_active, role, last_login_at, created_at, updated_at, refresh_token
		FROM users
		WHERE refresh_token = $1
	`
	err := r.db.Conn.QueryRowContext(ctx, query, refreshToken).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.Role,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.RefreshToken,
	)
	if err != nil {
		logger.Log.Error("Error getting user by refresh token", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info("User not found", zap.String("refresh_token", refreshToken))
			return User{}, ErrUserNotFound
		}
		logger.Log.Error("Error getting user by refresh token", zap.Error(err))
		return User{}, fmt.Errorf("failed to get user by refresh token: %w", err)
	}
	logger.Log.Info("Successfully got user by refresh token", zap.String("refresh_token", refreshToken))
	return user, nil
}

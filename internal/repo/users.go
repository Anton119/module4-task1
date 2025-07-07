package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
}

// UserRepository - обёртка для работы с пользователями в БД
type UserRepository struct {
	db *DB
}

// NewUserRepository - конструктор репозитория пользователей
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser добавляет нового пользователя в базу
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
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByUsername возвращает пользователя по username или ошибку ErrUserNotFound
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
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("failed to get user by username: %w", err)
	}
	return user, nil
}

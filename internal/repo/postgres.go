package repo

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"module4-task1/internal/logger"

	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sql.DB
}

func NewPostgres(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("Failed to connect to database", zap.String("dsn", dsn), zap.Error(err))
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Log.Error("Failed to ping database", zap.String("dsn", dsn), zap.Error(err))
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	logger.Log.Info("Successfully connected to database", zap.String("dsn", dsn))
	return &DB{Conn: db}, nil
}

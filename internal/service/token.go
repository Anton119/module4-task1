package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"module4-task1/internal/logger"
	"os"
	"time"
)

var secretKey []byte

func secretKeyInit() {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Info("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logger.Log.Info("Error loading JWT_SECRET")

	}
	logger.Log.Info("Successfully loaded JWT_SECRET")
	secretKey = []byte(secret)
}

func GenerateJWT(userID, username, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"name": username,
		"role": role,
		"exp":  time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshTokenStr() string {
	return uuid.NewString()
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Log.Error("Unexpected signing method")
			return nil, jwt.ErrTokenSignatureInvalid
		}
		logger.Log.Info("Successfully validated JWT")
		return secretKey, nil
	})

	if err != nil {
		logger.Log.Info("Error validating JWT", zap.Error(err))
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		logger.Log.Info("Successfully validated JWT", zap.Any("claims", claims))
		return claims, nil
	}

	logger.Log.Info("Successfully validated JWT", zap.Any("claims", token.Claims))
	return nil, jwt.ErrTokenInvalidClaims
}

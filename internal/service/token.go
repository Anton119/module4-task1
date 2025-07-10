package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var secretKey []byte

func secretKeyInit() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or failed to load")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in environment")
	}
	secretKey = []byte(secret)
}

func GenerateJWT(userID, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"name": username,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // токен живёт 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

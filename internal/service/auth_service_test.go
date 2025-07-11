package service

import (
	"testing"
	"time"
)

func TestAuthService_ValidateToken(t *testing.T) {
	secretKey = []byte("test-secret")
	service := &AuthService{}

	t.Run("valid token", func(t *testing.T) {
		validToken, err := GenerateJWT("user-id-123", "testuser", "user", 1*time.Hour)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		userID, err := service.ValidateToken(validToken)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if userID != "user-id-123" {
			t.Errorf("expected userID 'user-id-123', got '%s'", userID)
		}
	})

	t.Run("invalid token signature", func(t *testing.T) {
		badToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmYWtlLWlkIn0.invalidsignature"

		_, err := service.ValidateToken(badToken)
		if err == nil {
			t.Error("expected error for invalid token, got nil")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		expiredToken, err := GenerateJWT("user-id", "testuser", "user", -1*time.Hour)
		if err != nil {
			t.Fatalf("failed to generate expired token: %v", err)
		}

		_, err = service.ValidateToken(expiredToken)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

}

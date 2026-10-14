package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Config holds JWT settings (only Secret needed for validation).
type Config struct {
	Secret string
}

// JWTService handles JWT token validation.
type JWTService struct {
	secret string
}

// Claims represents JWT claims (must match auth-service structure!).
type Claims struct {
	UserID      uuid.UUID `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWTService.
func NewJWTService(config Config) *JWTService {
	return &JWTService{secret: config.Secret}
}

// ValidateToken parses and validates a JWT token.
func (s *JWTService) ValidateToken(tokenString string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

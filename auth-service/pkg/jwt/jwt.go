// pkg/jwt/jwt.go
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config holds JWT service settings.
type Config struct {
	Secret            string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secret            string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

// TokenPair contains access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// NewJWTService creates a new JWTService.
func NewJWTService(config Config) *JWTService {
	return &JWTService{
		secret:            config.Secret,
		accessExpiration:  config.AccessExpiration,
		refreshExpiration: config.RefreshExpiration,
	}
}

// GenerateTokenPair creates both access and refresh tokens.
func (s *JWTService) GenerateTokenPair(claims jwt.Claims) (*TokenPair, error) {
	accessToken, err := s.generateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpiration.Seconds()),
	}, nil
}

// ValidateToken parses and validates a JWT token string.
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

// GetAccessExpiration returns the access token expiration duration.
func (s *JWTService) GetAccessExpiration() time.Duration {
	return s.accessExpiration
}

// GetRefreshExpiration returns the refresh token expiration duration.
func (s *JWTService) GetRefreshExpiration() time.Duration {
	return s.refreshExpiration
}

// generateToken creates a signed JWT token.
func (s *JWTService) generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secret))
}

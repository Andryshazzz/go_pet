package usersservice

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secret            string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

// Claims represents the JWT claims for authenticated users.
type Claims struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

// TokenPair contains access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// NewJWTService creates a new JWTService.
func NewJWTService(secret string, accessExp, refreshExp time.Duration) *JWTService {
	return &JWTService{
		secret:            secret,
		accessExpiration:  accessExp,
		refreshExpiration: refreshExp,
	}
}

// GenerateTokenPair creates both access and refresh tokens for a user.
func (s *JWTService) GenerateTokenPair(userID, phoneNumber string) (*TokenPair, error) {
	accessToken, err := s.generateToken(userID, phoneNumber, s.accessExpiration)
	if err != nil {
		return nil, fmt.Errorf("Generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(userID, phoneNumber, s.refreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("Generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessExpiration.Seconds()),
	}, nil
}

// ValidateToken parses and validates a JWT token string.
// Returns the claims if the token is valid.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return claims, nil
}

// generateToken creates a signed JWT token.
func (s *JWTService) generateToken(userID, phoneNumber string, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID:      userID,
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secret))
}
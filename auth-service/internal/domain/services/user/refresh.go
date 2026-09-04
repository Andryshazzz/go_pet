package usersservice

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/internal/domain/entity"
	jwtpkg "github.com/Andryshazzz/go_pet/pkg/jwt"
)

// RefreshTokenRequest contains the refresh token for obtaining new tokens.
type RefreshTokenRequest struct {
	RefreshToken string
}

// RefreshTokenResponse contains the new token pair.
type RefreshTokenResponse struct {
	TokenPair *jwtpkg.TokenPair
}

// RefreshToken validates the refresh token and issues a new token pair.
func (s *UsersService) RefreshToken(
	ctx context.Context,
	req RefreshTokenRequest,
) (*RefreshTokenResponse, error) {
	claims := &entity.Claims{}
	
	err := s.jwtService.ValidateToken(req.RefreshToken, claims)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(claims)
	if err != nil {
		return nil, fmt.Errorf("generate new tokens: %w", err)
	}

	return &RefreshTokenResponse{
		TokenPair: tokenPair,
	}, nil
}

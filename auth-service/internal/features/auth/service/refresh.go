package usersservice

import (
	"context"
	"fmt"

	apperrors "github.com/Andryshazzz/go_pet/internal/core/errors"
)

// RefreshTokenRequest contains the refresh token for obtaining new tokens.
type RefreshTokenRequest struct {
	RefreshToken string
}

// RefreshTokenResponse contains the new token pair.
type RefreshTokenResponse struct {
	TokenPair *TokenPair
}

// RefreshToken validates the refresh token and issues a new token pair.
// The old refresh token becomes invalid (rotation for security).
func (s *UsersService) RefreshToken(
	ctx context.Context,
	req RefreshTokenRequest,
) (*RefreshTokenResponse, error) {
	claims, err := s.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid refresh token: %w: %w", err, apperrors.ErrUnauthorized)
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(claims.UserID, claims.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("Generate new tokens: %w", err)
	}

	return &RefreshTokenResponse{
		TokenPair: tokenPair,
	}, nil
}
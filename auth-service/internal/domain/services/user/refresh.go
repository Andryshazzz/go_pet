package usersservice

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/internal/domain/entity"
	apperrors "github.com/Andryshazzz/go_pet/pkg/errors"
	jwtpkg "github.com/Andryshazzz/go_pet/pkg/jwt"
)

// RefreshTokenRequest contains the refresh token for obtaining new tokens.
type RefreshTokenRequest struct {
	AccessToken  string
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
	accessClaims := &entity.Claims{}
	
	errAccessToken := s.jwtService.ValidateToken(req.AccessToken, accessClaims)
	if errAccessToken != nil {
		return nil, fmt.Errorf("invalid access token: %w", apperrors.ErrInvalidToken)
	}

	refreshClaims := &entity.Claims{}

	errRefreshToken := s.jwtService.ValidateToken(req.RefreshToken, refreshClaims)
	if errRefreshToken != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", apperrors.ErrInvalidToken)
	}

	if accessClaims.UserID != refreshClaims.UserID {
		return nil, fmt.Errorf("access and refresh token mismatch: %w", apperrors.ErrInvalidToken)
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("generate new tokens: %w", err)
	}

	return &RefreshTokenResponse{
		TokenPair: tokenPair,
	}, nil
}

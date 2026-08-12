package usersservice

import (
	"context"
	"fmt"

	apperrors "github.com/Andryshazzz/go_pet/internal/core/errors"
)

// LoginUserRequest contains the credentials needed for authentication.
type LoginUserRequest struct {
	PhoneNumber string
	Password string
}

// LoginUserResponse contains the result of successful authentication.
type LoginUserResponse struct {
	User UserResponse
	TokenPair *TokenPair
}

// UserResponse represents public user data returned to clients.
type UserResponse struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
}

// Login authenticates a user by phone number and password.
// It returns user data and a JWT token pair on success.
func (s *UsersService) Login(
	ctx context.Context,
	req LoginUserRequest,
) (*LoginUserResponse, error) {
	user, err := s.usersRepository.FindByPhone(ctx, req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("find user by phone: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user with phone %s not found: %w", req.PhoneNumber, apperrors.ErrNotFound)
	}

	if !checkPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid password: %w", apperrors.ErrInvalidArgument)
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &LoginUserResponse{
		User: UserResponse{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			FullName:    user.FullName,
		},
		TokenPair: tokenPair,
	}, nil
}
package usersservice

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/pkg/domain"
	apperrors "github.com/Andryshazzz/go_pet/pkg/errors"
)

// RegisterUserRequest contains the data needed for registration.
type RegisterUserRequest struct {
	PhoneNumber string
	Password    string
	FullName    string
}

// RegisterUserResponse contains the result of successful registration.
type RegisterUserResponse struct {
	User      domain.User
	TokenPair *TokenPair
}

// Register creates a new user after validation and uniqueness check.
// Returns the created user and a JWT token pair.
func (s *UsersService) Register(
	ctx context.Context,
	req RegisterUserRequest,
) (*RegisterUserResponse, error) {
	existingUser, err := s.usersRepository.FindByPhone(ctx, req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("Check phone existence: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("User with phone %s already exists: %w", req.PhoneNumber, apperrors.ErrUserAlreadyExists)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("Hash password: %w", err)
	}

	user := domain.NewUser(req.PhoneNumber, passwordHash, req.FullName)
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("Validate user: %w", err)
	}

	createdUser, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(createdUser.ID, createdUser.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("Generate tokens: %w", err)
	}

	return &RegisterUserResponse{
		User:      createdUser,
		TokenPair: tokenPair,
	}, nil
}

package usersservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/Andryshazzz/go_pet/internal/domain/entity"
	apperrors "github.com/Andryshazzz/go_pet/pkg/errors"
	jwt "github.com/Andryshazzz/go_pet/pkg/jwt"
)

// RegisterUserRequest contains the data needed for registration.
type RegisterUserRequest struct {
	PhoneNumber string
	Password    string
	FullName    string
}

// RegisterUserResponse contains the result of successful registration.
type RegisterUserResponse struct {
	User      entity.User
	TokenPair *jwt.TokenPair
}

// Register creates a new user after validation and uniqueness check.
// Returns the created user and a JWT token pair.
func (s *UsersService) Register(
	ctx context.Context,
	req RegisterUserRequest,
) (*RegisterUserResponse, error) {
	existingUser, err := s.usersRepository.FindByPhone(ctx, req.PhoneNumber)
	if err != nil && !errors.Is(err, apperrors.ErrNotFoundUser) {
		return nil, fmt.Errorf("check phone existence: %w", err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("user with phone %s already exists: %w", req.PhoneNumber, apperrors.ErrUserAlreadyExists)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := entity.NewUser(req.PhoneNumber, passwordHash, req.FullName)
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validate user: %w", err)
	}

	claims := entity.NewClaims(user.ID, user.PhoneNumber, s.jwtService.GetAccessExpiration())
	tokenPair, err := s.jwtService.GenerateTokenPair(claims)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	createdUser, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	
	return &RegisterUserResponse{
		User:      createdUser,
		TokenPair: tokenPair,
	}, nil
}

package usersservice

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/internal/core/domain"
)

// CreateUser validates the user and persists it using the repository.
// It returns the created user with server-generated fields.
func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.usersRepository.CreateUser(ctx, user)
	
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

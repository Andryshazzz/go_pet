package usersservice

import (
	"context"

	"github.com/Andryshazzz/go_pet/internal/core/domain"
)

// UsersRepository defines the interface for user persistence operations.
type UsersRepository interface {
	// CreateUser persists a new user and returns the created user with
    // any server-generated fields populated.
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

// UsersService contains business logic for user operations.
type UsersService struct {
	usersRepository UsersRepository
}

// NewUsersService creates a new UsersService with the given repository.
func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

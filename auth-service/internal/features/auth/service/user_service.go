package usersservice

import (
	"context"

	"github.com/Andryshazzz/go_pet/pkg/domain"
	jwt "github.com/Andryshazzz/go_pet/pkg/domain/jwt"
)

// UsersRepository defines the interface for user persistence operations.
type UsersRepository interface {
	// CreateUser persists a new user and returns the created user with
	// any server-generated fields populated.
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)

	// FindByPhone looks up a user by phone number.
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)
}

// UsersService contains business logic for user operations.
type UsersService struct {
	usersRepository UsersRepository
	jwtService      *jwt.JWTService
}

// NewUsersService creates a new UsersService with the given repository.
func NewUsersService(usersRepository UsersRepository, jwtService *jwt.JWTService) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		jwtService:      jwtService,
	}
}

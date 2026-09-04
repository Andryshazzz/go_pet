package usersservice

import (
	"context"

	"github.com/Andryshazzz/go_pet/internal/domain/entity"
	jwt "github.com/Andryshazzz/go_pet/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// UsersRepository defines the interface for user persistence operations.
type UsersRepository interface {
	// CreateUser persists a new user and returns the created user with
	// any server-generated fields populated.
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)

	// FindByPhone looks up a user by phone number.
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
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

// hashPassword creates a bcrypt hash of the password.
// Uses bcrypt.DefaultCost (10 rounds).
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

// checkPassword compares a password with its bcrypt hash.
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

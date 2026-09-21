package entity

import (
	"fmt"
	"regexp"

	apperrors "github.com/Andryshazzz/go_pet/pkg/errors"
	"github.com/google/uuid"
)

// User represents a user entity in the system.
//
// Fields:
//   - ID: unique identifier (UUID, empty for new users)
//   - PhoneNumber: phone number used as login (required, unique)
//   - PasswordHash: bcrypt hash of the password (never store plain text!)
//   - FullName: user's full name (2-100 characters)
type User struct {
	ID           uuid.UUID
	PhoneNumber  string
	PasswordHash string
	FullName     string
}

// NewUser creates a User with the given ID.
// Use this when loading an existing user from the database.
func NewUser(
	phoneNumber string,
	passwordHash string,
	fullName string,
) User {
	return User{
		ID:           uuid.New(),
		PhoneNumber:  phoneNumber,
		PasswordHash: passwordHash,
		FullName:     fullName,
	}
}

// Validate checks that the user data meets business rules.
// Returns an error if any field is invalid.
//
// Rules:
//   - FullName: 2-100 characters
//   - PhoneNumber (optional): 10-15 digits, optionally starting with +
func (u *User) Validate() error {
	phoneLen := len([]rune(u.PhoneNumber))
	if phoneLen < 10 || phoneLen > 15 {
		return fmt.Errorf(
			"Invalid PhoneNumber length %d (must be 10-15): %w",
			phoneLen,
			apperrors.ErrInvalidArgument,
		)
	}

	re := regexp.MustCompile(`^\+?[0-9]+$`)
	if !re.MatchString(u.PhoneNumber) {
		return fmt.Errorf(
			"Invalid PhoneNumber format %q (must be digits, optional +): %w",
			u.PhoneNumber,
			apperrors.ErrInvalidArgument,
		)
	}

	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 2 || fullNameLen > 100 {
		return fmt.Errorf(
			"Invalid FullName length %d (must be 2-100): %w",
			fullNameLen,
			apperrors.ErrInvalidArgument,
		)
	}

	return nil
}

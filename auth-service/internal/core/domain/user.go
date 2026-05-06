package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/Andryshazzz/go_pet/internal/core/errors"
)

type User struct {
	ID          int
	FullName    string
	PhoneNumber *string
}

func NewUser(
	id int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

}

func NewUserUninitialization(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializationID,
		fullName,
		phoneNumber,
	)
}

func (u *User) Validat() error {
	fullNameLenght := len([]rune(u.FullName))
	if fullNameLenght < 2 || fullNameLenght > 100 {
		return fmt.Errorf(
			"invalid `fullName` len: %d: %w",
			fullNameLenght,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		PhoneNumberLen := len([]rune(*u.PhoneNumber))
		if PhoneNumberLen < 10 || PhoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				fullNameLenght,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %d: %w",
				fullNameLenght,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

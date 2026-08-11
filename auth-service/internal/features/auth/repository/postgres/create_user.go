package usersrepository

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/internal/core/domain"
)

// CreateUser inserts a new user into the database and returns
// the created user with the generated ID.
func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO users (full_name, phone_number)
	VALUES ($1, $2)
	RETURNING id, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.FullName,
		user.PhoneNumber,
	)

	return userDomain, nil
}

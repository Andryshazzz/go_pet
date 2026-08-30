package usersrepository

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/pkg/domain"
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
		INSERT INTO users (phone_number, password_hash, full_name)
		VALUES ($1, $2, $3)
		RETURNING id, phone_number, password_hash, full_name;
	`

	row := r.pool.QueryRow(ctx, query,
		user.PhoneNumber,
		user.PasswordHash,
		user.FullName,
	)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.PhoneNumber,
		&userModel.PasswordHash,
		&userModel.FullName,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.NewUser(
		userModel.PhoneNumber,
		userModel.PasswordHash,
		userModel.FullName,
	), nil
}

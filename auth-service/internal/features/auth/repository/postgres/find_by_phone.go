package usersrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Andryshazzz/go_pet/pkg/domain"
	"github.com/jackc/pgx/v5"
)

// FindByPhone looks up a user by their phone number.
// Returns nil, nil if no user is found (not an error).
// Returns error only on database failures.
func (r *UsersRepository) FindByPhone(
	ctx context.Context,
	phone string,
) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, phone_number, password_hash, full_name
		FROM users
		WHERE phone_number = $1
	`

	var userModel UserModel
	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&userModel.ID,
		&userModel.PhoneNumber,
		&userModel.PasswordHash,
		&userModel.FullName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Find user by phone: %w", err)
	}

	user := domain.NewUser(
		userModel.PhoneNumber,
		userModel.PasswordHash,
		userModel.FullName,
	)
	return &user, nil
}

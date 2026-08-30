package usersrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Andryshazzz/go_pet/internal/domain/entity"
	postgrespool "github.com/Andryshazzz/go_pet/pkg/database/postgres/pool"
	apperrors "github.com/Andryshazzz/go_pet/pkg/errors"
	"github.com/jackc/pgx/v5"
)

// UsersRepository implements data access operations for users
// using PostgreSQL as the storage backend.
type UsersRepository struct {
	pool postgrespool.Pool
}

// NewUsersRepository creates a new UsersRepository with the given database pool.
// The pool must satisfy the postgrespool.Pool interface.
func NewUsersRepository(
	pool postgrespool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}

// CreateUser inserts a new user into the database and returns
// the created user with the generated ID.
func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user entity.User,
) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO users (phone_number, password_hash, full_name)
		VALUES ($1, $2, $3)
		RETURNING id, phone_number, password_hash, full_name;
	`

	row := r.pool.QueryRow(ctx, query,
		user.ID,
		user.PhoneNumber,
		user.PasswordHash,
		user.FullName,
	)

	var output entity.User

	err := row.Scan(
		&output.ID,
		&output.PhoneNumber,
		&output.PasswordHash,
		&output.FullName,
	)

	if err != nil {
		return entity.User{}, fmt.Errorf("scan error: %w", err)
	}

	return output, nil
}

// FindByPhone looks up a user by their phone number.
// Returns nil, nil if no user is found (not an error).
// Returns error only on database failures.
func (r *UsersRepository) FindByPhone(
	ctx context.Context,
	phone string,
) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, phone_number, password_hash, full_name
		FROM users
		WHERE phone_number = $1
	`

	var output entity.User

	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&output.ID,
		&output.PhoneNumber,
		&output.PasswordHash,
		&output.FullName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with phone %s: %w", phone, apperrors.ErrNotFoundUser)
		}
		return nil, fmt.Errorf("Find user by phone: %w", err)
	}
	
	return &output, nil
}

package usersrepository

import postgrespool "github.com/Andryshazzz/go_pet/pkg/database/postgres/pool"

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

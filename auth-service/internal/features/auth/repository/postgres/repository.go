package users_postgres_repository

import core_postgres_pool "github.com/Andryshazzz/go_pet/internal/core/repository/postgres/pool"

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}

package users_postrgres_repository

import (
	core_postgres_pool "github.com/cephalopagus/bkv-golang-todo/internal/core/repository/postgres/pool"
)

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(pool core_postgres_pool.Pool) *UsersRepository {
	return &UsersRepository{pool: pool}
}

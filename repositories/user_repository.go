package repositories

import (
	"database/sql"
)

type UserRepoOpt struct {
	Db *sql.DB
}

type UserRepository interface {
}

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(urOpt *UserRepoOpt) UserRepository {
	return &UserRepositoryPostgres{
		db: urOpt.Db,
	}
}

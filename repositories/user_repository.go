package repositories

import (
	"context"
	"database/sql"

	"github.com/Eggi19/simple-social-media/custom_errors"
	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories/queries"
)

type UserRepoOpt struct {
	Db *sql.DB
}

type UserRepository interface {
	RegisterUser(ctx context.Context, req entities.User) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(urOpt *UserRepoOpt) UserRepository {
	return &UserRepositoryPostgres{
		db: urOpt.Db,
	}
}

func (r *UserRepositoryPostgres) RegisterUser(ctx context.Context, req entities.User) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, queries.CreateUser)
	} else {
		stmt, err = r.db.PrepareContext(ctx, queries.CreateUser)
	}

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryPostgres) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	u := entities.User{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, queries.GetUserByEmail, email).Scan(&u.Id, &u.Name, &u.Email, &u.Password)
	} else {
		err = r.db.QueryRowContext(ctx, queries.GetUserByEmail, email).Scan(&u.Id, &u.Name, &u.Email, &u.Password)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound()
		}
		return nil, err
	}

	return &u, nil
}

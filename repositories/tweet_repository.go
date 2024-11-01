package repositories

import (
	"context"
	"database/sql"

	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories/queries"
)

type TweetRepoOpt struct {
	Db *sql.DB
}

type TweetRepository interface {
	CreateTweet(ctx context.Context, req entities.Tweet) error
}

type TweetRepositoryPostgres struct {
	db *sql.DB
}

func NewTweetRepositoryPostgres(trOpt *TweetRepoOpt) TweetRepository {
	return &TweetRepositoryPostgres{
		db: trOpt.Db,
	}
}

func (r *TweetRepositoryPostgres) CreateTweet(ctx context.Context, req entities.Tweet) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, queries.CreateTweet)
	} else {
		stmt, err = r.db.PrepareContext(ctx, queries.CreateTweet)
	}

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, req.Tweet, req.UserId)
	if err != nil {
		return err
	}

	return nil
}

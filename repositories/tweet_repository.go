package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"firebase.google.com/go/v4/db"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories/queries"
)

type TweetRepoOpt struct {
	Db              *sql.DB
	FirebseDbClient *db.Client
}

type TweetRepository interface {
	CreateTweet(ctx context.Context, req entities.Tweet) error
	UpdateTweetLikes(ctx context.Context, tweetId int64, likes int64) error
	GetTweetLikes(ctx context.Context, tweetId int64) (*dtos.TweetLikesResponse, error)
}

type TweetRepositoryDb struct {
	db              *sql.DB
	FirebseDbClient *db.Client
}

func NewTweetRepositoryDb(trOpt *TweetRepoOpt) TweetRepository {
	return &TweetRepositoryDb{
		db:              trOpt.Db,
		FirebseDbClient: trOpt.FirebseDbClient,
	}
}

func (r *TweetRepositoryDb) CreateTweet(ctx context.Context, req entities.Tweet) error {
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

func (r *TweetRepositoryDb) UpdateTweetLikes(ctx context.Context, tweetId int64, likes int64) error {
	reference := r.FirebseDbClient.NewRef(fmt.Sprintf(`tweets/%d`, tweetId))

	likesData := map[string]interface{}{
		"likes": likes,
	}

	err := reference.Set(ctx, likesData)
	if err != nil {
		return err
	}

	return nil
}

func (r *TweetRepositoryDb) GetTweetLikes(ctx context.Context, tweetId int64) (*dtos.TweetLikesResponse, error) {
	var likes *dtos.TweetLikesResponse
	reference := r.FirebseDbClient.NewRef(fmt.Sprintf(`tweets/%d`, tweetId))

	err := reference.Get(ctx, &likes)
	if err != nil {
		return nil, err
	}

	return likes, nil
}

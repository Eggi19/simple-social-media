package repositories

import (
	"context"
	"database/sql"

	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories/queries"
)

type CommentRepoOpt struct {
	Db *sql.DB
}

type CommentRepository interface {
	CreateComment(ctx context.Context, req entities.Comment) error
}

type CommentRepositoryPostgres struct {
	db *sql.DB
}

func NewCommentRepositoryPostgres(crOpt *CommentRepoOpt) CommentRepository {
	return &CommentRepositoryPostgres{
		db: crOpt.Db,
	}
}

func (r *CommentRepositoryPostgres) CreateComment(ctx context.Context, req entities.Comment) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, queries.CreateComment)
	} else {
		stmt, err = r.db.PrepareContext(ctx, queries.CreateComment)
	}

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, req.Comment, req.UserId, req.TweetId)
	if err != nil {
		return err
	}

	return nil
}

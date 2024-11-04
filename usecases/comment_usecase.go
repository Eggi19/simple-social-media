package usecases

import (
	"context"

	"firebase.google.com/go/messaging"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories"
)

type CommentUsecaseOpts struct {
	CommentRepository repositories.CommentRepository
	FirebaseClient    *messaging.Client
}

type CommentUsecase interface {
	CreateComment(ctx context.Context, userId int64, req dtos.CreateCommentRequest) error
}

type CommentUsecaseImpl struct {
	CommentRepository repositories.CommentRepository
	FirebaseClient    *messaging.Client
}

func NewCommentUsecaseImpl(tuOpts *CommentUsecaseOpts) CommentUsecase {
	return &CommentUsecaseImpl{
		CommentRepository: tuOpts.CommentRepository,
		FirebaseClient:    tuOpts.FirebaseClient,
	}
}

func (u *CommentUsecaseImpl) CreateComment(ctx context.Context, userId int64, req dtos.CreateCommentRequest) error {
	err := u.CommentRepository.CreateComment(ctx, entities.Comment{
		Comment: req.Comment,
		UserId:  userId,
		TweetId: req.TweetId,
	})
	if err != nil {
		return err
	}

	return nil
}

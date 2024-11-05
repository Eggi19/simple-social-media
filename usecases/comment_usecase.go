package usecases

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories"
)

type CommentUsecaseOpts struct {
	CommentRepository       repositories.CommentRepository
	UserRepository          repositories.UserRepository
	Transactor              repositories.Transactor
	FirebaseMessagingClient *messaging.Client
}

type CommentUsecase interface {
	CreateComment(ctx context.Context, userId int64, req dtos.CreateCommentRequest) error
}

type CommentUsecaseImpl struct {
	CommentRepository       repositories.CommentRepository
	UserRepository          repositories.UserRepository
	Transactor              repositories.Transactor
	FirebaseMessagingClient *messaging.Client
}

func NewCommentUsecaseImpl(tuOpts *CommentUsecaseOpts) CommentUsecase {
	return &CommentUsecaseImpl{
		CommentRepository:       tuOpts.CommentRepository,
		UserRepository:          tuOpts.UserRepository,
		Transactor:              tuOpts.Transactor,
		FirebaseMessagingClient: tuOpts.FirebaseMessagingClient,
	}
}

func (u *CommentUsecaseImpl) CreateComment(ctx context.Context, userId int64, req dtos.CreateCommentRequest) error {
	data, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.CommentRepository.CreateComment(ctx, entities.Comment{
			Comment: req.Comment,
			UserId:  userId,
			TweetId: req.TweetId,
		})
		if err != nil {
			return nil, err
		}

		user, err := u.UserRepository.GetUserIdByTweetId(ctx, req.TweetId)
		if err != nil {
			return nil, err
		}

		return user, nil
	})
	if err != nil {
		return err
	}

	user := data.(*entities.User)

	// set message payload.
	message := &messaging.Message{
		Token: user.FcmToken.String,
		Notification: &messaging.Notification{
			Title: fmt.Sprintf(`%s likes your tweet`, user.Name),
		},
	}

	_, err = u.FirebaseMessagingClient.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

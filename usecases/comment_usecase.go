package usecases

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories"
)

type CommentUsecaseOpts struct {
	CommentRepository repositories.CommentRepository
	UserRepository    repositories.UserRepository
	Transactor        repositories.Transactor
	FirebaseClient    *messaging.Client
}

type CommentUsecase interface {
	CreateComment(ctx context.Context, userId int64, req dtos.CreateCommentRequest) error
}

type CommentUsecaseImpl struct {
	CommentRepository repositories.CommentRepository
	UserRepository    repositories.UserRepository
	Transactor        repositories.Transactor
	FirebaseClient    *messaging.Client
}

func NewCommentUsecaseImpl(tuOpts *CommentUsecaseOpts) CommentUsecase {
	return &CommentUsecaseImpl{
		CommentRepository: tuOpts.CommentRepository,
		UserRepository:    tuOpts.UserRepository,
		Transactor:        tuOpts.Transactor,
		FirebaseClient:    tuOpts.FirebaseClient,
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
	
		userData, err := u.UserRepository.GetUserById(ctx, user.Id)
		if err != nil {
			return nil, err
		}

		return userData, nil
	})
	if err != nil {
		return err
	}

	userData := data.(*entities.User)

	// set message payload.
	message := &messaging.Message{
		Token: userData.FcmToken,
		Notification: &messaging.Notification{
			Title: fmt.Sprintf(`%s likes your tweet`, userData.Name),
		},
	}

	_, err = u.FirebaseClient.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

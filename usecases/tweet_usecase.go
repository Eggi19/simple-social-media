package usecases

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories"
)

type TweetUsecaseOpts struct {
	TweetRepository         repositories.TweetRepository
	UserRepository          repositories.UserRepository
	FirebaseMessagingClient *messaging.Client
}

type TweetUsecase interface {
	CreateTweet(ctx context.Context, userId int64, req dtos.CreateTweetRequest) error
	LikeTweet(ctx context.Context, tweetId int64) error
}

type TweetUsecaseImpl struct {
	TweetRepository         repositories.TweetRepository
	UserRepository          repositories.UserRepository
	FirebaseMessagingClient *messaging.Client
}

func NewTweetUsecaseImpl(tuOpts *TweetUsecaseOpts) TweetUsecase {
	return &TweetUsecaseImpl{
		TweetRepository:         tuOpts.TweetRepository,
		UserRepository:          tuOpts.UserRepository,
		FirebaseMessagingClient: tuOpts.FirebaseMessagingClient,
	}
}

func (u *TweetUsecaseImpl) CreateTweet(ctx context.Context, userId int64, req dtos.CreateTweetRequest) error {
	err := u.TweetRepository.CreateTweet(ctx, entities.Tweet{
		Tweet:  req.Tweet,
		UserId: userId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *TweetUsecaseImpl) LikeTweet(ctx context.Context, tweetId int64) error {
	var likes int64 = 1
	tweet, err := u.TweetRepository.GetTweetLikes(ctx, tweetId)
	if err != nil {
		return err
	}

	if tweet != nil {
		likes = tweet.Likes + 1
	}

	user, err := u.UserRepository.GetUserIdByTweetId(ctx, tweetId)
	if err != nil {
		return err
	}

	err = u.TweetRepository.UpdateTweetLikes(ctx, tweetId, likes)
	if err != nil {
		return err
	}

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

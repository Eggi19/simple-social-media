package usecases

import (
	"context"

	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/entities"
	"github.com/Eggi19/simple-social-media/repositories"
	"github.com/Eggi19/simple-social-media/utils"
)

type TweetUsecaseOpts struct {
	TweetRepository    repositories.TweetRepository
}

type TweetUsecase interface {
	CreateTweet(ctx context.Context, userId int64, req dtos.CreateTweetRequest) error
}

type TweetUsecaseImpl struct {
	HashAlgorithm     utils.Hasher
	AuthTokenProvider utils.AuthTokenProvider
	TweetRepository    repositories.TweetRepository
}

func NewTweetUsecaseImpl(tuOpts *TweetUsecaseOpts) TweetUsecase {
	return &TweetUsecaseImpl{
		TweetRepository:    tuOpts.TweetRepository,
	}
}

func (u *TweetUsecaseImpl) CreateTweet(ctx context.Context, userId int64, req dtos.CreateTweetRequest) error {
	err := u.TweetRepository.CreateTweet(ctx, entities.Tweet{
		Tweet: req.Tweet,
		UserId: userId,
	})
	if err != nil {
		return err
	}

	return nil
}

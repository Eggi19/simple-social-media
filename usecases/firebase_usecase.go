package usecases

import (
	"context"

	"github.com/Eggi19/simple-social-media/dtos"
	"github.com/Eggi19/simple-social-media/repositories"
)

type FirebaseUsecaseOpts struct {
	FirebaseRepository repositories.FirebaseRepository
}

type FirebaseUsecase interface {
	SubscribeTopic(ctx context.Context, req dtos.SubscribeTopicRequest) error
	UnsubscribeTopic(ctx context.Context, req dtos.UnsubscribeTopicRequest) error
}

type FirebaseUsecaseImpl struct {
	FirebaseRepository repositories.FirebaseRepository
}

func NewFirebaseUsecaseImpl(fuOpts *FirebaseUsecaseOpts) FirebaseUsecase {
	return &FirebaseUsecaseImpl{
		FirebaseRepository: fuOpts.FirebaseRepository,
	}
}

func (u *FirebaseUsecaseImpl) SubscribeTopic(ctx context.Context, req dtos.SubscribeTopicRequest) error {
	fcmTokens := []string{
		req.FcmToken,
	}

	err := u.FirebaseRepository.SubsribeTopic(ctx, fcmTokens, req.Topic)
	if err != nil {
		return err
	}

	return nil
}

func (u *FirebaseUsecaseImpl) UnsubscribeTopic(ctx context.Context, req dtos.UnsubscribeTopicRequest) error {
	fcmTokens := []string{
		req.FcmToken,
	}

	err := u.FirebaseRepository.UnsubsribeTopic(ctx, fcmTokens, req.Topic)
	if err != nil {
		return err
	}

	return nil
}

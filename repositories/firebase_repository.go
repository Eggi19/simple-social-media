package repositories

import (
	"context"

	"firebase.google.com/go/messaging"
)

type FirebaseRepoOpt struct {
	FirebaseClient *messaging.Client
}

type FirebaseRepository interface {
	SubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error
	UnsubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error
}

type FirebaseClientRepository struct {
	FirebaseClient *messaging.Client
}

func NewFirebaseRepositoryPostgres(frOpt *FirebaseRepoOpt) FirebaseRepository {
	return &FirebaseClientRepository{
		FirebaseClient: frOpt.FirebaseClient,
	}
}

func (r *FirebaseClientRepository) SubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error {
	_, err := r.FirebaseClient.SubscribeToTopic(ctx, fcmTokens, topic)
	if err != nil {
		return err
	}

	return nil
}

func (r *FirebaseClientRepository) UnsubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error {
	_, err := r.FirebaseClient.UnsubscribeFromTopic(ctx, fcmTokens, topic)
	if err != nil {
		return err
	}

	return nil
}

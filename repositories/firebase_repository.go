package repositories

import (
	"context"

	"firebase.google.com/go/v4/messaging"
)

type FirebaseRepoOpt struct {
	FirebaseMessagingClient *messaging.Client
}

type FirebaseRepository interface {
	SubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error
	UnsubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error
}

type FirebaseMessagingClientRepository struct {
	FirebaseMessagingClient *messaging.Client
}

func NewFirebaseRepositoryPostgres(frOpt *FirebaseRepoOpt) FirebaseRepository {
	return &FirebaseMessagingClientRepository{
		FirebaseMessagingClient: frOpt.FirebaseMessagingClient,
	}
}

func (r *FirebaseMessagingClientRepository) SubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error {
	_, err := r.FirebaseMessagingClient.SubscribeToTopic(ctx, fcmTokens, topic)
	if err != nil {
		return err
	}

	return nil
}

func (r *FirebaseMessagingClientRepository) UnsubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error {
	_, err := r.FirebaseMessagingClient.UnsubscribeFromTopic(ctx, fcmTokens, topic)
	if err != nil {
		return err
	}

	return nil
}

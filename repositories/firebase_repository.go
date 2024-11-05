package repositories

import (
	"context"

	"firebase.google.com/go/db"
	"firebase.google.com/go/messaging"
)

type FirebaseRepoOpt struct {
	FirebaseMessagingClient *messaging.Client
	FirebaseDbClient        *db.Client
}

type FirebaseRepository interface {
	SubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error
	UnsubsribeTopic(ctx context.Context, fcmTokens []string, topic string) error
	WriteData(ctx context.Context, ref string, data map[string]interface{}) error
}

type FirebaseMessagingClientRepository struct {
	FirebaseMessagingClient *messaging.Client
	FirebaseDbClient        *db.Client
}

func NewFirebaseRepositoryPostgres(frOpt *FirebaseRepoOpt) FirebaseRepository {
	return &FirebaseMessagingClientRepository{
		FirebaseMessagingClient: frOpt.FirebaseMessagingClient,
		FirebaseDbClient:        frOpt.FirebaseDbClient,
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

func (r *FirebaseMessagingClientRepository) WriteData(ctx context.Context, ref string, data map[string]interface{}) error {
	reference := r.FirebaseDbClient.NewRef(ref)

	err := reference.Set(ctx, data); 
	if err != nil {
        return err
    }

	return nil
}

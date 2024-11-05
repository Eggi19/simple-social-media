package config

import (
	"context"

	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	MessagingClient *messaging.Client
	DbClient        *db.Client
}

func NewFirebaseRepository(serviceAccountKeyPath string) (*FirebaseClient, error) {
	ctx := context.Background()

	// Initialize Firebase App
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		return nil, err
	}

	// Initialize FCM Client
	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	dbClient, err := app.Database(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseClient{
		MessagingClient: messagingClient,
		DbClient:        dbClient,
	}, nil
}

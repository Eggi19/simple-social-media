package config

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	MessagingClient *messaging.Client
	DbClient        *db.Client
}

func NewFirebaseRepository(config Config, serviceAccountKeyPath string) (*FirebaseClient, error) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{
        DatabaseURL: config.FirebaseDbUrl,
        ProjectID: config.FirebaseProjectId,
        StorageBucket: config.FirebaseStorageBucket,
    }

	// Initialize Firebase App
	app, err := firebase.NewApp(ctx, firebaseConfig, option.WithCredentialsFile(serviceAccountKeyPath))
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
		fmt.Println(err)
		return nil, err
	}

	return &FirebaseClient{
		MessagingClient: messagingClient,
		DbClient:        dbClient,
	}, nil
}

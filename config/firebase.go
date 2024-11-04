package config

import (
	"context"


    "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)


type FirebaseClient struct {
    Client *messaging.Client
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

    return &FirebaseClient{Client: messagingClient}, nil
}

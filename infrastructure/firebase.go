package infrastructure

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// NewFBApp creates new firebase app instance
func NewFBApp(logger Logger) *firebase.App {
	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		logger.Zap.Panic("Unable to load serviceAccountKey.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Zap.Fatalf("Firebase NewApp: %v", err)
	}

	logger.Zap.Info("Firebase app initialized 🎉🎉🎉")
	return app
}

// NewFBAuth creates new firebase auth client
func NewFBAuth(logger Logger, app *firebase.App) *auth.Client {

	ctx := context.Background()

	firebaseAuth, err := app.Auth(ctx)
	if err != nil {
		logger.Zap.Fatalf("Firebase Authentication: %v", err)
	}

	return firebaseAuth
}

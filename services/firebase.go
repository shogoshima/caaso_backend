package services

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	AuthClient *auth.Client
)

func InitFirebase() error {
	opt := option.WithCredentialsFile("./credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	// Initialize Auth client once
	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		return err
	}

	return nil
}

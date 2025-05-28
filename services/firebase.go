package services

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	AuthClient *auth.Client
)

func InitFirebase() error {
	opt := option.WithCredentialsFile("./credentials.json")
	config := &firebase.Config{ProjectID: "caaso-app"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Initialize Auth client once
	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		return err
	}

	return nil
}

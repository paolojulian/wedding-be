package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
)

var FirebaseApp *firebase.App
var FirestoreClient *firestore.Client

func InitFirebase() error {
	ctx := context.Background()
	serviceAccount := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}
	FirebaseApp = app

	return nil
}

func InitFirestore() error {
	ctx := context.Background()
	client, err := FirebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error getting firebase client: %v\n", err)
	}
	FirestoreClient = client

	return nil
}

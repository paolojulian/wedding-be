package testUtils

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadEnvFile(path string) {
	// Load the .env file
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
func SetupTestDB(t *testing.T) (*mongo.Database, func()) {
	// Load the .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		t.Fatal("MONGODB_URI is not set")
	}

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Get the collectoin
	mongoDB := client.Database("test_db")

	// Clean up after each test
	cleanup := func() {
		client.Disconnect(ctx)
		mongoDB.Collection("invitations").Drop(ctx)
	}

	return mongoDB, cleanup
}

package invitations

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/paolojulian/wedding-be/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateInvitationMessage(t *testing.T) {
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
	defer client.Disconnect(ctx)

	// Get the collectoin
	mongoDB := client.Database("test_db")
	service := NewInvitationMessageService(mongoDB)
	collection := mongoDB.Collection("invitation_message")

	// Clean up after each test
	defer collection.Drop(ctx)

	t.Run("Success Case", func(t *testing.T) {
		// Insert an invitation message first
		initialInvitationMessage := models.InvitationMessage{
			ID:      "qwe123",
			Message: "Test Message",
		}

		_, err := collection.InsertOne(ctx, initialInvitationMessage)
		if err != nil {
			t.Fatalf("Failed to insert test document: %v", err)
			return
		}

		newMessage := "expected message"
		updateErr := service.UpdateInvitationMessage(ctx, newMessage)

		if updateErr != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify the update
		var updatedInvitation models.InvitationMessage
		err = collection.FindOne(ctx, bson.M{}).Decode(&updatedInvitation)
		if err != nil {
			t.Fatalf("Failed to fetch updated document: %v", err)
		}

		if updatedInvitation.Message != newMessage {
			t.Errorf("Expected message to be %q, got %q", newMessage, updatedInvitation.Message)
		}
	})
}

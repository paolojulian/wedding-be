package invitations

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/paolojulian/wedding-be/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateInvitation(t *testing.T) {
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
	service := NewInvitationService(mongoDB)

	// Clean up after each test
	defer mongoDB.Collection("invitations").Drop(ctx)

	t.Run("Success Case", func(t *testing.T) {
		invitation := models.Invitation{
			Name:        "Test User",
			VoucherCode: "TEST123",
		}

		result, err := service.CreateInvitation(ctx, invitation)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result.ID == "" {
			t.Errorf("Expected ID to be generated")
		}
		if result.Status != "pending" {
			t.Errorf("Expected status to be 'pending', got %s", result.Status)
		}
		if result.Index != 1 {
			t.Errorf("Expected index to be 1, got %d", result.Index)
		}
		if result.Name != invitation.Name {
			t.Errorf("Expected name to be %s, got %s", invitation.Name, result.Name)
		}
		if result.VoucherCode != invitation.VoucherCode {
			t.Errorf("Expected voucher code to be %s, got %s", invitation.VoucherCode, result.VoucherCode)
		}
	})
}

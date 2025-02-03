package invitations

import (
	"context"
	"testing"

	"github.com/paolojulian/wedding-be/internal/models"
	testUtils "github.com/paolojulian/wedding-be/test-utils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateInvitation(t *testing.T) {
	mongoDB, cleanup := testUtils.SetupTestDB(t)
	service := NewInvitationService(mongoDB)

	defer cleanup()

	t.Run("Success Case", func(t *testing.T) {
		invitation := models.Invitation{
			Name:        "Test User",
			VoucherCode: "TEST123",
		}

		result, err := service.CreateInvitation(context.Background(), invitation)

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

func TestRespondToInvitation(t *testing.T) {
	ctx := context.Background()
	mongoDB, cleanup := testUtils.SetupTestDB(t)
	collection := mongoDB.Collection("invitations")
	service := NewInvitationService(mongoDB)

	defer cleanup()

	t.Run("Success Case", func(t *testing.T) {
		voucherCode := "qwe123"
		// We insert a test data first
		initialInvitation := models.Invitation{
			Name:          "Initial User",
			VoucherCode:   voucherCode,
			Index:         0,
			Status:        "pending",
			GuestsAllowed: 1,
			GuestsToBring: 0,
		}
		_, err := collection.InsertOne(ctx, initialInvitation)
		if err != nil {
			t.Fatalf("Failed to insert test document: %v", err)
		}

		respondRequest := RespondToInvitationRequest{
			Status:        "going",
			GuestsToBring: 1,
		}
		respondErr := service.RespondToInvitation(ctx, voucherCode, respondRequest)
		if respondErr != nil {
			t.Fatalf("Failed to respond to invitation: %v", respondErr)
		}

		// Verify the update
		var updatedInvitation models.Invitation
		err = collection.FindOne(ctx, bson.M{"voucher_code": voucherCode}).Decode(&updatedInvitation)
		if err != nil {
			t.Fatalf("Failed to fetch updated document: %v", err)
		}

		expectedInvitation := models.Invitation{
			ID:            updatedInvitation.ID,
			Name:          "Initial User",
			VoucherCode:   voucherCode,
			Index:         0,
			Status:        respondRequest.Status,
			GuestsAllowed: 1,
			GuestsToBring: respondRequest.GuestsToBring,
		}

		if updatedInvitation != expectedInvitation {
			t.Errorf("Expected status to be %q, got %q", expectedInvitation, updatedInvitation)
		}
	})
}

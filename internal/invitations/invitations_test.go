package invitations

import (
	"context"
	"testing"

	"github.com/paolojulian/wedding-be/internal/models"
	testUtils "github.com/paolojulian/wedding-be/test-utils"
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

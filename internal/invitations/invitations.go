package invitations

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paolojulian/wedding-be/internal/firebase"
	"github.com/paolojulian/wedding-be/internal/models"
	"google.golang.org/api/iterator"
)

func GetList(c *gin.Context) {
	iter := firebase.FirestoreClient.Collection("invitations").Documents(context.Background())
	var invitations []models.Invitation

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		data := doc.Data()
		log.Printf("Document data: %v", data)
		// Print individual fields for debugging
		log.Printf("Voucher Code: %v", data["voucher_code"])
		log.Printf("Name: %v", data["name"])
		log.Printf("Status: %v", data["status"])
		log.Printf("Guests Allowed: %v", data["guests_allowed"])
		log.Printf("Guests To Bring: %v", data["guests_to_bring"])

		invitation := models.Invitation{
			ID:            doc.Ref.ID,
			VoucherCode:   data["voucher_code"].(string),
			Name:          data["name"].(string),
			Status:        data["status"].(string),
			GuestsAllowed: int(data["guests_allowed"].(int64)),
			GuestsToBring: int(data["guests_to_bring"].(int64)),
		}
		log.Printf("Mapped Invitation: %+v", invitation) // Log the mapped invitation
		invitations = append(invitations, invitation)
	}

	c.JSON(http.StatusOK, invitations)
}
func RespondToInvitation() {}
func CreateInvitation()    {}
func EditInvitation()      {}

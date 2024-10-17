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

		invitation := models.Invitation{
			ID:            doc.Ref.ID,
			VoucherCode:   data["voucher_code"].(string),
			Name:          data["name"].(string),
			Status:        data["status"].(string),
			GuestsAllowed: int(data["guests_allowed"].(int64)),
			GuestsToBring: int(data["guests_to_bring"].(int64)),
		}
		invitations = append(invitations, invitation)
	}

	c.JSON(http.StatusOK, invitations)
}
func RespondToInvitation(c *gin.Context) {
	var newInvitation models.Invitation

	if err := c.BindJSON(&newInvitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
}
func CreateInvitation(c *gin.Context) {
	var newInvitation models.Invitation

	if err := c.BindJSON(&newInvitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request, check your form values"})
		return
	}

	_, _, err := firebase.FirestoreClient.Collection("invitations").Add(context.Background(), map[string]interface{}{
		"id":              newInvitation.ID,
		"index":           newInvitation.Index,
		"voucher_code":    newInvitation.VoucherCode,
		"name":            newInvitation.Name,
		"status":          newInvitation.Status,
		"guests_allowed":  newInvitation.GuestsAllowed,
		"guests_to_bring": newInvitation.GuestsToBring,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating invitation"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newInvitation)
}
func EditInvitation() {}

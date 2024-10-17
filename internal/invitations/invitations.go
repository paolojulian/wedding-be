package invitations

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paolojulian/wedding-be/internal/firebase"
	"github.com/paolojulian/wedding-be/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvitationService struct {
	collection *mongo.Collection
}

func NewInvitationService(db *mongo.Database) *InvitationService {
	return &InvitationService{
		collection: db.Collection("invitations"),
	}
}

func (s *InvitationService) GetList(c context.Context) ([]models.Invitation, error) {
	var invitations []models.Invitation
	iter, err := s.collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	defer iter.Close(c)

	if err := iter.All(c, &invitations); err != nil {
		return nil, err
	}

	return invitations, nil
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

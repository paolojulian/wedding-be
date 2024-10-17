package invitations

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paolojulian/wedding-be/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *InvitationService) CreateInvitation(c context.Context, invitation models.Invitation) (*models.Invitation, error) {
	// Validate required fields
	if invitation.Name == "" {
		return nil, ErrNameIsRequired
	}

	// Voucher code will be generated on the front-end side
	if invitation.VoucherCode == "" {
		return nil, ErrVoucherCodeIsRequired
	}

	// Set default values
	invitation.Status = "pending" // default status
	invitation.Index = 1

	result, err := s.collection.InsertOne(c, invitation)
	if err != nil {
		return nil, err
	}

	invitation.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return &invitation, nil
}

func (s *InvitationService) DeleteInvitation(c context.Context, ID string) error {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return ErrInvalidIDFormat
	}

	result, err := s.collection.DeleteOne(c, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return nil
	}

	return nil
}

func RespondToInvitation(c *gin.Context) {
	var newInvitation models.Invitation

	if err := c.BindJSON(&newInvitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
}
func EditInvitation() {}

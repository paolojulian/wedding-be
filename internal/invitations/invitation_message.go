package invitations

import (
	"context"

	"github.com/paolojulian/wedding-be/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/paolojulian/wedding-be/pkg/db"
)

type InvitationMessageService struct {
	collection *mongo.Collection
}

type UpdateInvitationMessageRequest struct {
	Message *string `json:"message"`
}

func NewInvitationMessageService(mongoDB *mongo.Database) *InvitationMessageService {
	return &InvitationMessageService{
		collection: mongoDB.Collection(db.InvitationMessageCollection),
	}
}

func (s *InvitationMessageService) GetInvitationMessage(c context.Context) (models.InvitationMessage, error) {
	var invitationMessage models.InvitationMessage
	err := s.collection.FindOne(c, bson.D{}).Decode(&invitationMessage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.InvitationMessage{}, ErrInvitationMessageNotFound
		}
		return models.InvitationMessage{}, err
	}

	return invitationMessage, nil

}

func (s *InvitationMessageService) UpdateInvitationMessage(c context.Context, invitationMessage string) error {
	updateDoc := bson.M{
		"message": invitationMessage,
	}

	var existingInvitationMessage models.InvitationMessage
	err := s.collection.FindOne(c, bson.D{}).Decode(&existingInvitationMessage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrInvitationMessageNotFound
		}
		return err
	}

	result, err := s.collection.UpdateOne(c, bson.M{"_id": existingInvitationMessage.ID}, bson.M{"$set": updateDoc})
	if err != nil {
		return ErrCannotUpdateInDB
	}

	if result.MatchedCount == 0 {
		return ErrInvitationMessageNotFound
	}

	return nil
}

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
			return models.InvitationMessage{}, ErrInvitationNotFound
		}
		return models.InvitationMessage{}, err
	}

	return invitationMessage, nil

}

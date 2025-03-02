package invitations

import (
	"context"

	app_config "github.com/paolojulian/wedding-be/internal/config"
	"github.com/paolojulian/wedding-be/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/paolojulian/wedding-be/pkg/db"
)

type UpdateInvitationRequest struct {
	Name          *string `json:"name"` // Using pointers to detect if field was provided
	VoucherCode   *string `json:"voucher_code"`
	Status        *string `json:"status"`
	GuestsAllowed *int    `json:"guests_allowed"` // Pointer allows detecting explicit 0
	GuestsToBring *int    `json:"guests_to_bring"`
	Index         *int    `json:"index"`
}

type RespondToInvitationRequest struct {
	Status        string `json:"status" bindings:"required"`
	GuestsToBring int    `json:"guests_to_bring" bindings:"required"`
}

type InvitationService struct {
	collection       *mongo.Collection
	appConfigService *app_config.AppConfigService
}

func NewInvitationService(mongoDB *mongo.Database) *InvitationService {
	return &InvitationService{
		collection:       mongoDB.Collection(db.InvitationsCollection),
		appConfigService: app_config.NewAppConfigService(mongoDB),
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
	invitation.GuestsToBring = 0

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

func (s *InvitationService) UpdateInvitation(c context.Context, ID string, invitation UpdateInvitationRequest) error {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return ErrInvalidIDFormat
	}

	updateDoc := bson.M{}

	if invitation.Name != nil {
		updateDoc["name"] = invitation.Name
	}

	if invitation.Status != nil {
		updateDoc["status"] = invitation.Status
	}

	if invitation.VoucherCode != nil {
		updateDoc["voucher_code"] = invitation.VoucherCode
	}

	if invitation.GuestsAllowed != nil {
		updateDoc["guests_allowed"] = invitation.GuestsAllowed
	}

	if invitation.GuestsToBring != nil {
		updateDoc["guests_to_bring"] = invitation.GuestsToBring
	}

	if len(updateDoc) == 0 {
		return ErrNoFieldsToUpdate
	}

	result, err := s.collection.UpdateOne(c, bson.M{"_id": objID}, bson.M{"$set": updateDoc})
	if err != nil {
		return ErrCannotUpdateInDB
	}

	if result.MatchedCount == 0 {
		return ErrInvitationNotFound
	}

	return nil
}

func (s *InvitationService) RespondToInvitation(c context.Context, VoucherCode string, respondReq RespondToInvitationRequest) error {
	isLocked, err := s.appConfigService.GetIsLocked(c)
	if err != nil {
		return err
	}

	if isLocked {
		return ErrIsAlreadyLocked
	}

	filter := bson.M{"voucher_code": VoucherCode}
	updateDoc := bson.M{
		"status":          respondReq.Status,
		"guests_to_bring": respondReq.GuestsToBring,
	}

	result, err := s.collection.UpdateOne(c, filter, bson.M{"$set": updateDoc})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrInvitationNotFound
	}

	return nil
}

func (s *InvitationService) GetInvitationByVoucherCode(c context.Context, VoucherCode string) (models.Invitation, error) {
	filter := bson.M{"voucher_code": VoucherCode}

	var invitation models.Invitation
	err := s.collection.FindOne(c, filter).Decode(&invitation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Invitation{}, ErrInvitationNotFound
		}
		return models.Invitation{}, err
	}

	return invitation, nil
}

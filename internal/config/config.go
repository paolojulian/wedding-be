package app_config

import (
	"context"

	"github.com/paolojulian/wedding-be/internal/models"
	"github.com/paolojulian/wedding-be/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfigService struct {
	collection *mongo.Collection
}

func NewAppConfigService(mongoDB *mongo.Database) *AppConfigService {
	return &AppConfigService{
		collection: mongoDB.Collection(db.ConfigCollection),
	}
}

/**
 * Get the config for isLocked (to determine if the rsvp is already closed)
 */
func (s *AppConfigService) GetIsLocked(c context.Context) (bool, error) {
	filter := bson.M{"name": ConfigNameIsLocked}

	var config models.Config
	err := s.collection.FindOne(c, filter).Decode(&config)
	if err != nil {
		return false, err
	}

	return config.Value == "1", nil
}

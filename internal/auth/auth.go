package auth

import (
	"context"

	"github.com/paolojulian/wedding-be/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	collection *mongo.Collection
}

func NewAuthService(db *mongo.Database) *AuthService {
	return &AuthService{
		collection: db.Collection("users"),
	}
}

func (s *AuthService) Login(c context.Context, username, password string) (string, error) {
	var user models.User
	err := s.collection.FindOne(c, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", ErrUserNotFound
	}

	// Verify password
	if !checkPasswordHash(password, user.Password) {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := generateJWT(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

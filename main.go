package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	// Internals

	"github.com/paolojulian/wedding-be/internal/auth"
	"github.com/paolojulian/wedding-be/internal/models"
	"github.com/paolojulian/wedding-be/internal/utils"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var invitations = []models.Invitation{
	{
		ID:            "1",
		VoucherCode:   "123456",
		Name:          "John Doe",
		Status:        "going",
		GuestsAllowed: 2,
		GuestsToBring: models.IntPtr(1),
	},
	{
		ID:            "2",
		VoucherCode:   "223456",
		Name:          "Paolo Vincent Julian",
		Status:        "pending",
		GuestsAllowed: 1,
		GuestsToBring: nil,
	},
}

func main() {

	router := gin.Default()

	// Initialize Firebase Service account
	ctx := context.Background()
	serviceAccount := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
		return
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error getting firebase client: %v\n", err)
		return
	}
	defer client.Close()

	// Allow CORS from localhost
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/invitations", utils.AuthMiddleware(), getInvitations)
	router.POST("/invitations", utils.AuthMiddleware(), postInvitation)
	router.PUT("/invitations/:id", utils.AuthMiddleware(), editInvitation)
	router.PUT("/invitations/respond/:voucherCode", utils.AuthMiddleware(), respondToInvitation)

	// Authentication endpoints
	router.GET("/me", utils.AuthMiddleware(), auth.ValidateLoggedInUser)
	router.POST("/login", auth.Login)
	router.POST("/logout", auth.Logout)

	router.Run("localhost:8080")
}

func getInvitations(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, invitations)
}

func postInvitation(c *gin.Context) {
	var newInvitation models.Invitation

	if err := c.BindJSON(&newInvitation); err != nil {
		return
	}

	invitations = append(invitations, newInvitation)
	c.IndentedJSON(http.StatusCreated, newInvitation)
}

func editInvitation(c *gin.Context) {
	id := c.Param("id")

	// If found, update the invitation
	var updatedInvitation models.Invitation
	if err := c.BindJSON(&updatedInvitation); err != nil {
		return
	}

	// Update the invitation details
	for i, invitation := range invitations {
		if invitation.ID == id {
			invitations[i] = updatedInvitation
			break
		}
	}

	// If not found, return 404
	if updatedInvitation.ID == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Invitation not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedInvitation)
}

func respondToInvitation(c *gin.Context) {
	voucherCode := c.Param("voucherCode")

	var updatedInvitation models.Invitation
	if err := c.BindJSON(&updatedInvitation); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	for i, invitation := range invitations {
		if invitation.VoucherCode == voucherCode {
			invitations[i].GuestsToBring = updatedInvitation.GuestsToBring
			invitations[i].Status = updatedInvitation.Status
			break
		}
	}

	if updatedInvitation.VoucherCode == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Invitation not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedInvitation)
}

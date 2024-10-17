package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// Internals

	"github.com/paolojulian/wedding-be/internal/auth"
	"github.com/paolojulian/wedding-be/internal/firebase"
	"github.com/paolojulian/wedding-be/internal/invitations"
	"github.com/paolojulian/wedding-be/internal/models"
	"github.com/paolojulian/wedding-be/internal/utils"
	"github.com/paolojulian/wedding-be/pkg/db"
)

var invitationArr = []models.Invitation{
	{
		ID:            "1",
		VoucherCode:   "123456",
		Name:          "John Doe",
		Status:        "going",
		GuestsAllowed: 2,
		GuestsToBring: 0,
	},
	{
		ID:            "2",
		VoucherCode:   "223456",
		Name:          "Paolo Vincent Julian",
		Status:        "pending",
		GuestsAllowed: 1,
		GuestsToBring: 0,
	},
}

func main() {
	router := gin.Default()

	client := db.ConnectMongoDB()
	defer db.DisconnectMongoDB()

	// Initialize the services
	authService := auth.NewAuthService(client.Database("wedding_db"))
	authHandler := auth.NewHandler(authService)

	firebase.InitFirebase()
	firebase.InitFirestore()

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

	// Invitation endpoints
	router.GET("/invitations", utils.AuthMiddleware(), invitations.GetList)
	router.GET("/test/invitations", invitations.GetList)
	router.POST("/invitations", utils.AuthMiddleware(), invitations.CreateInvitation)
	router.PUT("/invitations/:id", utils.AuthMiddleware(), editInvitation)
	router.PUT("/invitations/respond/:voucherCode", utils.AuthMiddleware(), respondToInvitation)

	// Authentication endpoints
	router.GET("/me", utils.AuthMiddleware(), authHandler.ValidateLoggedInUser)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	router.Run("0.0.0.0:8080")
}

func editInvitation(c *gin.Context) {
	id := c.Param("id")

	// If found, update the invitation
	var updatedInvitation models.Invitation
	if err := c.BindJSON(&updatedInvitation); err != nil {
		return
	}

	// Update the invitation details
	for i, invitation := range invitationArr {
		if invitation.ID == id {
			invitationArr[i] = updatedInvitation
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

	for i, invitation := range invitationArr {
		if invitation.VoucherCode == voucherCode {
			invitationArr[i].GuestsToBring = updatedInvitation.GuestsToBring
			invitationArr[i].Status = updatedInvitation.Status
			break
		}
	}

	if updatedInvitation.VoucherCode == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Invitation not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedInvitation)
}

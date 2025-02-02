package main

import (
	"github.com/gin-gonic/gin"

	// Internals

	"github.com/paolojulian/wedding-be/config"
	"github.com/paolojulian/wedding-be/internal/auth"
	"github.com/paolojulian/wedding-be/internal/invitations"
	"github.com/paolojulian/wedding-be/pkg/db"
)

func main() {
	router := gin.Default()

	client := db.ConnectMongoDB()
	defer db.DisconnectMongoDB()

	// Initialize the services
	authService := auth.NewAuthService(client.Database(db.DatabaseName))
	authHandler := auth.NewHandler(authService)
	invitationService := invitations.NewInvitationService(client.Database(db.DatabaseName))
	invitationMessageService := invitations.NewInvitationMessageService(client.Database(db.DatabaseName))
	invitationHandler := invitations.NewHandler(invitationService, invitationMessageService)

	// Get URIs from environment variable
	appURI := config.GetAppURI()
	adminURI := config.GetAdminURI()

	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed
		if origin == adminURI || origin == appURI {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Set-Cookie")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
			// Increase max age for better caching of CORS headers
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Protected endpoints
	router.GET("/invitations", auth.AuthMiddleware(), invitationHandler.GetList)
	router.POST("/invitations", auth.AuthMiddleware(), invitationHandler.CreateInvitation)
	router.PUT("/invitations/:id", auth.AuthMiddleware(), invitationHandler.UpdateInvitation)
	router.DELETE("/invitations/:id", auth.AuthMiddleware(), invitationHandler.DeleteInvitation)
	router.GET("/invitation-message", auth.AuthMiddleware(), invitationHandler.GetInvitationMessage)

	// Public endpoints
	router.PUT("/invitations/respond/:voucher_code", invitationHandler.RespondToInvitation)
	router.GET("/invitations/respond/:voucher_code", invitationHandler.GetInvitationForRespond)

	// Authentication endpoints
	router.GET("/me", auth.AuthMiddleware(), authHandler.ValidateLoggedInUser)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	router.Run("0.0.0.0:8080")
}

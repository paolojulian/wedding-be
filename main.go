package main

import (
	"github.com/gin-gonic/gin"

	// Internals

	"github.com/paolojulian/wedding-be/internal/auth"
	"github.com/paolojulian/wedding-be/internal/firebase"
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
	invitationHandler := invitations.NewHandler(invitationService)

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
	router.GET("/invitations", auth.AuthMiddleware(), invitationHandler.GetList)
	router.GET("/test/invitations", invitationHandler.GetList)
	router.POST("/invitations", auth.AuthMiddleware(), invitationHandler.CreateInvitation)
	router.PUT("/invitations/:id", auth.AuthMiddleware(), invitationHandler.UpdateInvitation)
	router.DELETE("/invitations/:id", auth.AuthMiddleware(), invitationHandler.DeleteInvitation)
	router.PUT("/invitations/respond/:voucher_code", auth.AuthMiddleware(), invitationHandler.RespondToInvitation)
	router.GET("/invitations/respond/:voucher_code", auth.AuthMiddleware(), invitationHandler.GetInvitationForRespond)

	// Authentication endpoints
	router.GET("/me", auth.AuthMiddleware(), authHandler.ValidateLoggedInUser)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	router.Run("0.0.0.0:8080")
}

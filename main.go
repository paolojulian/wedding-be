package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

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

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Get URIs from environment variable
	appURI := os.Getenv("APP_URI")
	if appURI == "" {
		log.Fatal("APP_URI environment variable not set")
	}
	adminURI := os.Getenv("ADMIN_URI")
	if adminURI == "" {
		log.Fatal("ADMIN_URI environment variable not set")
	}

	// Allow CORS from localhost
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == adminURI || origin == appURI {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Protected endpoints
	router.GET("/invitations", auth.AuthMiddleware(), invitationHandler.GetList)
	router.GET("/test/invitations", invitationHandler.GetList)
	router.POST("/invitations", auth.AuthMiddleware(), invitationHandler.CreateInvitation)
	router.PUT("/invitations/:id", auth.AuthMiddleware(), invitationHandler.UpdateInvitation)
	router.DELETE("/invitations/:id", auth.AuthMiddleware(), invitationHandler.DeleteInvitation)

	// Public endpoints
	router.PUT("/invitations/respond/:voucher_code", invitationHandler.RespondToInvitation)
	router.GET("/invitations/respond/:voucher_code", invitationHandler.GetInvitationForRespond)

	// Authentication endpoints
	router.GET("/me", auth.AuthMiddleware(), authHandler.ValidateLoggedInUser)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	router.Run("0.0.0.0:8080")
}

package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paolojulian/wedding-be/internal/firebase"
)

func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Validate the credentials
	usersCollection := firebase.FirestoreClient.Collection("users")
	doc, err := usersCollection.Doc("1").Get(context.Background())
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	var user struct {
		Password string `firestore:"password"`
	}

	if err := doc.DataTo(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error processing user data"})
		return
	}

	if user.Password != credentials.Password {
		log.Default().Println("Password entered: ", credentials, "Password in DB: ", user.Password)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	setCookieHandler(c.Writer, "authTokenSampleqwerty")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func Logout(c *gin.Context) {
	// Clear the auth token cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func ValidateLoggedInUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "You are logged in"})
}

func setCookieHandler(w http.ResponseWriter, cookieValue string) {
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}

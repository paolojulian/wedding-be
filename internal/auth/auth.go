package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	if credentials.Username != "admin" || credentials.Password != "qwe123" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
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

package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *AuthService
}

func NewHandler(service *AuthService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Validate the credentials
	token, err := h.Service.Login(c, credentials.Username, credentials.Password)
	if err != nil {
		if err == ErrUserNotFound || err == ErrInvalidCredentials {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	setCookieHandler(c.Writer, token)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func (h *Handler) Logout(c *gin.Context) {
	// Clear the auth token cookie
	c.SetCookie(authCookieName, "", -1, "/", "", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) ValidateLoggedInUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "You are logged in"})
}

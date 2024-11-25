package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
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

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Login successfully.",
		"token":   token,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) ValidateLoggedInUser(c *gin.Context) {
	// The auth middleware has already validated it, so we don't need to add extra functions here
	c.IndentedJSON(http.StatusOK, gin.H{"message": "You are logged in"})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return

		}

		// Extract the token (assuming format: "Bearer <token>")
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Malformed token"})
			c.Abort()
			return
		}

		// Validate the JWT token
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure that the token method used is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the claims in the context for use in downstream handlers
		c.Set("claims", claims)
		c.Next()
	}
}

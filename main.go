package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type invitation struct {
	ID            string `json:"id"`
	VoucherCode   string `json:"voucher_code"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	GuestsAllowed int    `json:"guests_allowed"`
	GuestsToBring *int   `json:"guests_to_bring"`
}

var invitations = []invitation{
	{
		ID:            "1",
		VoucherCode:   "123456",
		Name:          "John Doe",
		Status:        "going",
		GuestsAllowed: 2,
		GuestsToBring: intPtr(1),
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

func intPtr(i int) *int {
	return &i
}

func main() {
	router := gin.Default()
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

	router.GET("/invitations", authMiddleware(), getInvitations)
	router.POST("/invitations", authMiddleware(), postInvitation)
	router.PUT("/invitations/:id", authMiddleware(), editInvitation)
	router.PUT("/invitations/respond/:voucherCode", authMiddleware(), respondToInvitation)

	// Authentication endpoints
	router.GET("/me", authMiddleware(), checkMe)
	router.POST("/login", login)
	router.POST("/logout", logout)

	router.Run("localhost:8080")
}

func getInvitations(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, invitations)
}

func postInvitation(c *gin.Context) {
	var newInvitation invitation

	if err := c.BindJSON(&newInvitation); err != nil {
		return
	}

	invitations = append(invitations, newInvitation)
	c.IndentedJSON(http.StatusCreated, newInvitation)
}

func editInvitation(c *gin.Context) {
	id := c.Param("id")

	// If found, update the invitation
	var updatedInvitation invitation
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

	var updatedInvitation invitation
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

func login(c *gin.Context) {
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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Request.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		// TODO: Validate cookie
		c.Next()
	}
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

func logout(c *gin.Context) {
	// Clear the auth token cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func checkMe(c *gin.Context) {
	// If the user is logged in, return the user details
	c.IndentedJSON(http.StatusOK, gin.H{"username": "admin"})
}

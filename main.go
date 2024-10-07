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
	GuestsToBring int    `json:"guests_to_bring"`
}

var invitations = []invitation{
	{
		ID:            "1",
		VoucherCode:   "123456",
		Name:          "John Doe",
		Status:        "going",
		GuestsAllowed: 2,
		GuestsToBring: 1,
	},
}

func main() {
	router := gin.Default()
	router.GET("/invitations", getInvitations)
	router.POST("/invitations", postInvitation)
	router.PUT("/invitations/:id", editInvitation)
	router.PUT("/invitations/respond/:voucherCode", respondToInvitation)

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

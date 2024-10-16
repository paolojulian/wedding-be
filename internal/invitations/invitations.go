package invitations

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paolojulian/wedding-be/internal/database"
	"github.com/paolojulian/wedding-be/internal/models"
)

func GetList(c *gin.Context) {
	invitations, err := database.ReadInvitations()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, invitations)
}
func RespondToInvitation(c *gin.Context) {
	var newInvitation models.Invitation

	if err := c.BindJSON(&newInvitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
}
func CreateInvitation() {}
func EditInvitation()   {}

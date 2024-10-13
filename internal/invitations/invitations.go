package invitations

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paolojulian/wedding-be/internal/database"
)

func GetList(c *gin.Context) {
	invitations, err := database.ReadInvitations()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, invitations)
}
func RespondToInvitation() {}
func CreateInvitation()    {}
func EditInvitation()      {}

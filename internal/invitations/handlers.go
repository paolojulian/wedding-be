package invitations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvitationHandler struct {
	Service *InvitationService
}

func NewHandler(service *InvitationService) *InvitationHandler {
	return &InvitationHandler{
		Service: service,
	}
}

func (h *InvitationHandler) GetList(c *gin.Context) {
	invitations, err := h.Service.GetList(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, invitations)
}

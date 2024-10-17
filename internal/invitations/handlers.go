package invitations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paolojulian/wedding-be/internal/models"
)

type Handler struct {
	Service *InvitationService
}

func NewHandler(service *InvitationService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetList(c *gin.Context) {
	invitations, err := h.Service.GetList(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, invitations)
}

func (h *Handler) CreateInvitation(c *gin.Context) {
	type CreateInvitationRequest struct {
		Name          string `json:"name" binding:"required"`
		VoucherCode   string `json:"voucher_code,omitempty" binding:"required"`
		GuestsAllowed int    `json:"guests_allowed,omitempty" binding:"required"`
	}

	// Parse request body
	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Create invitation model from request
	invitation := models.Invitation{
		Name:          req.Name,
		VoucherCode:   req.VoucherCode,
		GuestsAllowed: req.GuestsAllowed,
		GuestsToBring: 0, // Initialize as 0
	}

	result, err := h.Service.CreateInvitation(c, invitation)
	if err != nil {
		if err == ErrNameIsRequired || err == ErrVoucherCodeIsRequired {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Incomplete form"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) DeleteInvitation(c *gin.Context) {
	ID := c.Param("id")

	err := h.Service.DeleteInvitation(c, ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (s *Handler) UpdateInvitation(c *gin.Context) {
	ID := c.Param("id")

	var req UpdateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	err := s.Service.UpdateInvitation(c, ID, req)
	if err != nil {
		switch err {
		case ErrInvitationNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invitation not found"})
		case ErrInvalidIDFormat:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		case ErrNoFieldsToUpdate:
			c.JSON(http.StatusOK, req)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to update the invitation"})
		}
		return
	}

	c.JSON(http.StatusOK, req)
}

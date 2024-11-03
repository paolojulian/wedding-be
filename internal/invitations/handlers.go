package invitations

import (
	"log"
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
		GuestsAllowed int    `json:"guests_allowed,omitempty"`
	}

	// Parse request body
	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error %v", err)
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
			c.JSON(http.StatusBadRequest, gin.H{"message": "No fields to update"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to update the invitation"})
		}
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) RespondToInvitation(c *gin.Context) {
	VoucherCode := c.Param("voucher_code")

	var req RespondToInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}

	err := h.Service.RespondToInvitation(c, VoucherCode, req)
	if err != nil {
		switch err {
		case ErrInvitationNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "Invitation Not Found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Responded!"})
}

func (h *Handler) GetInvitationForRespond(c *gin.Context) {
	VoucherCode := c.Param("voucher_code")

	invitation, err := h.Service.GetInvitationByVoucherCode(c, VoucherCode)
	if err != nil {
		switch err {
		case ErrInvitationNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "Invitation not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
	}

	c.JSON(http.StatusOK, invitation)
}

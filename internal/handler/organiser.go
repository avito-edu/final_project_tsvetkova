package handler

import (
	"net/http"
	"strconv"
	"swim_service/internal/dto"

	"github.com/gin-gonic/gin"
)

type OrganizerService interface {
	CreateOrganizer(organizerDTO dto.CreateOrganizerRequest) (*dto.OrganizerResponse, error)
	GetOrganizer(id int) (*dto.OrganizerResponse, error)
}

type OrganizerHandler struct {
	organizerService OrganizerService
}

func NewOrganizerHandler(organizerService OrganizerService) *OrganizerHandler {
	return &OrganizerHandler{
		organizerService: organizerService,
	}
}

// CreateOrganizer godoc
// @Summary Create a new organizer
// @Description Create a new competition organizer with the provided information
// @Tags organizers
// @Accept json
// @Produce json
// @Param organizer body dto.CreateOrganizerRequest true "Organizer information"
// @Success 201 {object} dto.OrganizerResponse "Successfully created organizer"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /organizers/create [post]
func (h *OrganizerHandler) CreateOrganizer(c *gin.Context) {
	var req dto.CreateOrganizerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.organizerService.CreateOrganizer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetOrganizer godoc
// @Summary Get organizer by ID
// @Description Retrieve organizer information by organizer ID
// @Tags organizers
// @Accept json
// @Produce json
// @Param id path integer true "Organizer ID"
// @Success 200 {object} dto.OrganizerResponse "Successfully retrieved organizer"
// @Failure 400 {object} map[string]string "Invalid organizer ID format"
// @Failure 404 {object} map[string]string "Organizer not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /organizers/get/{id} [get]
func (h *OrganizerHandler) GetOrganizer(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organizer ID"})
		return
	}

	resp, err := h.organizerService.GetOrganizer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organizer not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
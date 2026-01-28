package handler

import (
	"net/http"
	"strconv"
	"swim_service/internal/dto"

	"github.com/gin-gonic/gin"
)

type AthleteService interface {
	CreateAthlete(athleteDTO dto.CreateAthleteRequest) (*dto.AthleteResponse, error)
	GetAthlete(id int) (*dto.AthleteResponse, error)
	UpdateAthleteStatus(id int, status string) (*dto.AthleteResponse, error)
}

type AthleteHandler struct {
	athleteService AthleteService
}

func NewAthleteHandler(athleteService AthleteService) *AthleteHandler {
	return &AthleteHandler{
		athleteService: athleteService,
	}
}

// CreateAthlete godoc
// @Summary Create a new athlete
// @Description Create a new athlete with the provided information
// @Tags athletes
// @Accept json
// @Produce json
// @Param athlete body dto.CreateAthleteRequest true "Athlete information"
// @Success 201 {object} dto.AthleteResponse "Successfully created athlete"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /athletes/create [post]
func (h *AthleteHandler) CreateAthlete(c *gin.Context) {
	var req dto.CreateAthleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.athleteService.CreateAthlete(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetAthlete godoc
// @Summary Get athlete by ID
// @Description Retrieve athlete information by athlete ID
// @Tags athletes
// @Accept json
// @Produce json
// @Param id path integer true "Athlete ID"
// @Success 200 {object} dto.AthleteResponse "Successfully retrieved athlete"
// @Failure 400 {object} map[string]string "Invalid athlete ID format"
// @Failure 404 {object} map[string]string "Athlete not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /athletes/get/{id} [get]
func (h *AthleteHandler) GetAthlete(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid athlete ID"})
		return
	}

	resp, err := h.athleteService.GetAthlete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Athlete not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateStatus godoc
// @Summary Update athlete status
// @Description Update the status of an existing athlete
// @Tags athletes
// @Accept json
// @Produce json
// @Param id path integer true "Athlete ID"
// @Param status body dto.UpdateAthleteStatusRequest true "New status information"
// @Success 200 {object} dto.AthleteResponse "Successfully updated athlete status"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Athlete not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /athletes/update/{id} [patch]
func (h *AthleteHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid athlete ID"})
		return
	}

	var req dto.UpdateAthleteStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.athleteService.UpdateAthleteStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
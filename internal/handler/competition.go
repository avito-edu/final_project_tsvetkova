package handler

import (
	"net/http"
	"strconv"
	"swim_service/internal/dto"

	"github.com/gin-gonic/gin"
)

type CompetitionService interface {
	CreateCompetition(compDTO dto.CreateCompetitionRequest) (*dto.CompetitionResponse, error)
	GetCompetition(id int) (*dto.CompetitionResponse, error)
	GetAllCompetitions() ([]dto.CompetitionResponse, error)
	AddResult(compID int, resultDTO dto.AddResultRequest) (*dto.ResultResponse, error)
	GetAthleteProgress(athleteID int, distance string) ([]dto.ResultResponse, error)
	GetAllResultsByDistance(distance string) ([]dto.ResultResponse, error)
}

type CompetitionHandler struct {
	competitionService CompetitionService
}

func NewCompetitionHandler(competitionService CompetitionService) *CompetitionHandler {
	return &CompetitionHandler{
		competitionService: competitionService,
	}
}

// CreateCompetition godoc
// @Summary Create a new competition
// @Description Create a new swimming competition with the provided information
// @Tags competitions
// @Accept json
// @Produce json
// @Param competition body dto.CreateCompetitionRequest true "Competition information"
// @Success 201 {object} dto.CompetitionResponse "Successfully created competition"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /competitions/create [post]
func (h *CompetitionHandler) CreateCompetition(c *gin.Context) {
	var req dto.CreateCompetitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.competitionService.CreateCompetition(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetAllCompetitions godoc
// @Summary Get all competitions
// @Description Retrieve a list of all competitions
// @Tags competitions
// @Accept json
// @Produce json
// @Success 200 {array} dto.CompetitionResponse "Successfully retrieved competitions list"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /competitions/get-all [get]
func (h *CompetitionHandler) GetAllCompetitions(c *gin.Context) {
	resp, err := h.competitionService.GetAllCompetitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCompetition godoc
// @Summary Get competition by ID
// @Description Retrieve competition information by competition ID
// @Tags competitions
// @Accept json
// @Produce json
// @Param id path integer true "Competition ID"
// @Success 200 {object} dto.CompetitionResponse "Successfully retrieved competition"
// @Failure 400 {object} map[string]string "Invalid competition ID format"
// @Failure 404 {object} map[string]string "Competition not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /competitions/get/{id} [get]
func (h *CompetitionHandler) GetCompetition(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition ID"})
		return
	}

	resp, err := h.competitionService.GetCompetition(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Competition not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddResult godoc
// @Summary Add result to competition
// @Description Add a new result (swimmer's time) to an existing competition
// @Tags results
// @Accept json
// @Produce json
// @Param id path integer true "Competition ID"
// @Param result body dto.AddResultRequest true "Result information"
// @Success 201 {object} dto.ResultResponse "Successfully added result"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Competition not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /competitions/add/{id} [post]
func (h *CompetitionHandler) AddResult(c *gin.Context) {
	idStr := c.Param("id")

	compID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition ID"})
		return
	}

	var req dto.AddResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.competitionService.AddResult(compID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
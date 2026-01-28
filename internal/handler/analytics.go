package handler

import (
	"net/http"
	"strconv"
	"swim_service/internal/dto"
	"swim_service/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetAthleteAnalytics godoc
// @Summary Get athlete analytics
// @Description Get detailed analytics for a specific athlete including performance metrics and trends
// @Tags analytics
// @Accept json
// @Produce json
// @Param id path integer true "Athlete ID"
// @Param request body dto.AthleteAnalyticsRequest true "Analytics period filters"
// @Success 200 {object} dto.AthleteAnalyticsResponse "Successful response with athlete analytics"
// @Failure 400 {object} dto.ErrorResponse "Invalid athlete ID or request format"
// @Failure 404 {object} dto.ErrorResponse "Athlete not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /analytics/get/{id} [post]
func (h *AnalyticsHandler) GetAthleteAnalytics(c *gin.Context) {
	athleteIDStr := c.Param("id")
	
	athleteID, err := strconv.Atoi(athleteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid athlete ID format"})
		return
	}

	var req dto.AthleteAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	req.AthleteID = athleteID

	if req.PeriodStart.IsZero() {
		req.PeriodStart = time.Now().AddDate(0, -1, 0)
	}
	if req.PeriodEnd.IsZero() {
		req.PeriodEnd = time.Now()
	}

	analytics, err := h.analyticsService.GetAthleteAnalytics(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
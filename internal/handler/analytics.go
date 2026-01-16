package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"swim_service/internal/dto"
	"swim_service/internal/service"
	"time"

	"github.com/gorilla/mux"
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
// @Failure 400 {object} map[string]string "Invalid athlete ID or request format"
// @Failure 404 {object} map[string]string "Athlete not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /analytics/get/{id} [post]
func (h *AnalyticsHandler) GetAthleteAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	athleteIDStr := vars["id"]

	athleteID, err := strconv.Atoi(athleteIDStr)
	if err != nil {
		http.Error(w, "Invalid athlete ID", http.StatusBadRequest)
		return
	}

	var req dto.AthleteAnalyticsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.AthleteID = athleteID

	if req.PeriodStart.IsZero() {
		req.PeriodStart = time.Now().AddDate(0, -1, 0)
	}
	if req.PeriodEnd.IsZero() {
		req.PeriodEnd = time.Now()
	}

	analytics, err := h.analyticsService.GetAthleteAnalytics(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		http.Error(w, "Failed to make response", http.StatusInternalServerError)
	}
}

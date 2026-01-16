package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"swim_service/internal/dto"

	"github.com/gorilla/mux"
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
func (h *AthleteHandler) CreateAthlete(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAthleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.athleteService.CreateAthlete(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to make response", http.StatusInternalServerError)
	}
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
func (h *AthleteHandler) GetAthlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid athlete ID", http.StatusBadRequest)
		return
	}

	resp, err := h.athleteService.GetAthlete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp == nil {
		http.Error(w, "Athlete not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to make response", http.StatusInternalServerError)
	}
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
func (h *AthleteHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid athlete ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateAthleteStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.athleteService.UpdateAthleteStatus(id, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to make response", http.StatusInternalServerError)
	}
}

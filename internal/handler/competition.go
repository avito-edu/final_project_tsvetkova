package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"swim_service/internal/dto"

	"github.com/gorilla/mux"
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
func (h *CompetitionHandler) CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCompetitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.competitionService.CreateCompetition(req)
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

// GetAllCompetitions godoc
// @Summary Get all competitions
// @Description Retrieve a list of all competitions
// @Tags competitions
// @Accept json
// @Produce json
// @Success 200 {array} dto.CompetitionResponse "Successfully retrieved competitions list"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /competitions/get-all [get]
func (h *CompetitionHandler) GetAllCompetitions(w http.ResponseWriter, r *http.Request) {
	resp, err := h.competitionService.GetAllCompetitions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to make response", http.StatusInternalServerError)
	}
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
func (h *CompetitionHandler) GetCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid competition ID", http.StatusBadRequest)
		return
	}

	resp, err := h.competitionService.GetCompetition(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp == nil {
		http.Error(w, "Competition not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to make response", http.StatusInternalServerError)
	}
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
func (h *CompetitionHandler) AddResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	compID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid competition ID", http.StatusBadRequest)
		return
	}

	var req dto.AddResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.competitionService.AddResult(compID, req)
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

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"swim_service/internal/dto"

	"github.com/gorilla/mux"
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
func (h *OrganizerHandler) CreateOrganizer(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrganizerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.organizerService.CreateOrganizer(req)
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
func (h *OrganizerHandler) GetOrganizer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid organizer ID", http.StatusBadRequest)
		return
	}

	resp, err := h.organizerService.GetOrganizer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp == nil {
		http.Error(w, "Organizer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to make response", http.StatusInternalServerError)
	}
}

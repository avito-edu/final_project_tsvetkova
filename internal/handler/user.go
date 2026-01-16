package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"swim_service/internal/domain"
	"swim_service/internal/dto"
)

type UserService interface {
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
	RegisterSpecial(ctx context.Context, req dto.RegisterSpecialRequest) (*domain.User, error)
	RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*domain.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with login and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterUserRequest true "User registration data"
// @Success 200 {object} map[string]any "User successfully registered"
// @Failure 400 {object} map[string]string "Invalid JSON or request format"
// @Failure 409 {object} map[string]string "User already exists"
// @Router /register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userService.RegisterUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":    user.ID,
		"login": user.Login,
		"role":  user.Role,
	})
}

// RegisterSpecial godoc
// @Summary Register a special user (e.g., admin, organizer)
// @Description Register a user with special role and permissions
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterSpecialRequest true "Special user registration data"
// @Success 200 {object} map[string]any "Special user successfully registered"
// @Failure 400 {object} map[string]string "Invalid JSON or request format"
// @Failure 409 {object} map[string]string "User already exists"
// @Router /register-special [post]
func (h *UserHandler) RegisterSpecial(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterSpecialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userService.RegisterSpecial(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":    user.ID,
		"login": user.Login,
		"role":  user.Role,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and receive JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "User login credentials"
// @Success 200 {object} map[string]string "Login successful, returns JWT token"
// @Failure 400 {object} map[string]string "Invalid JSON or request format"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := h.userService.Login(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

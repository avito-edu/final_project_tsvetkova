package handler

import (
	"context"
	"net/http"
	"swim_service/internal/domain"
	"swim_service/internal/dto"

	"github.com/gin-gonic/gin"
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
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	user, err := h.userService.RegisterUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
func (h *UserHandler) RegisterSpecial(c *gin.Context) {
	var req dto.RegisterSpecialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	user, err := h.userService.RegisterSpecial(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	token, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
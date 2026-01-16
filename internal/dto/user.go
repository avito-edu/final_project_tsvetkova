package dto

// RegisterUserRequest represents user registration request
// @Description Request for registering a new regular user
type RegisterUserRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterSpecialRequest represents special user registration
// @Description Request for registering a user with special role (admin, organizer, etc.)
type RegisterSpecialRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=user admin"`
}

// LoginRequest represents user login request
// @Description Request for user authentication
type LoginRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

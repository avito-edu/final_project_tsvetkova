package dto

import "time"

// CreateAthleteRequest represents request for creating a new athlete
// @Description Request for creating a new athlete
type CreateAthleteRequest struct {
	Name      string     `json:"name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Status    string     `json:"status"`
	Gender    string     `json:"gender"`
}

// UpdateAthleteStatusRequest represents request for updating athlete status
// @Description Request for updating athlete status
type UpdateAthleteStatusRequest struct {
	Status string `json:"status"`
}

// AthleteResponse represents athlete information
// @Description Response with athlete details
type AthleteResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Status    string     `json:"status"`
	Gender    string     `json:"gender"`
}

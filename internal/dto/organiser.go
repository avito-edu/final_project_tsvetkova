package dto

// CreateOrganizerRequest represents request for creating an organizer
// @Description Request for creating a new competition organizer
type CreateOrganizerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// OrganizerResponse represents organizer information
// @Description Response with organizer details
type OrganizerResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

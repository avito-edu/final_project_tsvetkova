package dto

import "time"

// CreateCompetitionRequest represents request for creating a competition
// @Description Request for creating a new swimming competition
type CreateCompetitionRequest struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	OrganizerID int       `json:"organizer_id"`
}

// CompetitionResponse represents competition information
// @Description Response with competition details
type CompetitionResponse struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Date      time.Time        `json:"date"`
	Organizer int              `json:"organizer"`
	Results   []ResultResponse `json:"results"`
}

// AddResultRequest represents request for adding a result
// @Description Request for adding a swimmer's result to a competition
type AddResultRequest struct {
	AthleteID int     `json:"athlete_id"`
	Distance  string  `json:"distance"`
	TimeSec   float64 `json:"time_sec"`
	Place     int     `json:"place"`
}

// ResultResponse represents swimming result information
// @Description Response with swimmer's result details
type ResultResponse struct {
	AthleteID int     `json:"athlete_id"`
	Distance  string  `json:"distance"`
	TimeSec   float64 `json:"time_sec"`
	Place     int     `json:"place"`
}

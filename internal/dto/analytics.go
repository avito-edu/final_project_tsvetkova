package dto

import "time"

// AthleteAnalyticsRequest represents request for athlete analytics
// @Description Request for retrieving athlete analytics with optional period filters
type AthleteAnalyticsRequest struct {
	AthleteID   int       `json:"athlete_id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

// AthleteAnalyticsResponse represents response with athlete analytics
// @Description Detailed analytics data for an athlete
type AthleteAnalyticsResponse struct {
	AthleteID         int              `json:"athlete_id"`
	PeriodStart       time.Time        `json:"period_start"`
	PeriodEnd         time.Time        `json:"period_end"`
	TotalDistance     string           `json:"total_distance"`
	TotalCompetitions int              `json:"total_competitions"`
	AverageTime       float64          `json:"average_time"`
	BestTime          float64          `json:"best_time"`
	WorstTime         float64          `json:"worst_time"`
	Results           []ResultResponse `json:"results"`
	DistanceBreakdown []DistanceStats  `json:"distance_breakdown"`
}

// DistanceStats represents statistics for a specific swimming distance
// @Description Statistical data for swimming results by distance
type DistanceStats struct {
	Distance    string  `json:"distance"`
	Count       int     `json:"count"`
	AverageTime float64 `json:"average_time"`
	BestTime    float64 `json:"best_time"`
	WorstTime   float64 `json:"worst_time"`
}

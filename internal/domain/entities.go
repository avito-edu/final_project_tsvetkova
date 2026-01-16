package domain

import "time"

type Role string

const (
	RoleAthlet Role = "athlete"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string
	Role     Role `json:"role"`
}

type Athlete struct {
	ID        string
	Name      string
	BirthDate *time.Time
	Status    string
	Gender    string
}

type Organizer struct {
	ID    int
	Name  string
	Email string
}

type Competition struct {
	ID        int
	Name      string
	Date      time.Time
	Organizer int
	Results   []Result
}

type Result struct {
	AthleteID int
	Distance  string
	TimeSec   float64
	Place     int
}

type AthleteAnalytics struct {
	AthleteID         int
	PeriodStart       time.Time
	PeriodEnd         time.Time
	TotalDistance     string
	TotalCompetitions int
	AverageTime       float64
	BestTime          float64
	WorstTime         float64
	Results           []Result
}

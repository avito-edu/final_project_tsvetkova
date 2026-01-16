package handler

import (
	"net/http"
	"swim_service/internal/middleware"

	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
)

func NewRouter(
	athleteHandler *AthleteHandler,
	competitionHandler *CompetitionHandler,
	organizerHandler *OrganizerHandler,
	analyticsHandler *AnalyticsHandler,
	userHandler *UserHandler,
	enforcer *casbin.Enforcer,
) http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/athletes/create", athleteHandler.CreateAthlete).Methods("POST")
	mux.HandleFunc("/athletes/get/{id}", athleteHandler.GetAthlete).Methods("GET")
	mux.HandleFunc("/athletes/update/{id}", athleteHandler.UpdateStatus).Methods("PATCH")

	mux.HandleFunc("/competitions/create", competitionHandler.CreateCompetition).Methods("POST")
	mux.HandleFunc("/competitions/get/{id}", competitionHandler.GetCompetition).Methods("GET")
	mux.HandleFunc("/competitions/get-all", competitionHandler.GetAllCompetitions).Methods("GET")
	mux.HandleFunc("/competitions/add/{id}", competitionHandler.AddResult).Methods("POST")

	mux.HandleFunc("/organizers/create", organizerHandler.CreateOrganizer).Methods("POST")
	mux.HandleFunc("/organizers/get/{id}", organizerHandler.GetOrganizer).Methods("GET")

	mux.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")
	mux.HandleFunc("/register-special", userHandler.RegisterSpecial).Methods("POST")
	mux.HandleFunc("/login", userHandler.Login).Methods("POST")

	mux.HandleFunc("/analytics/get/{id}", analyticsHandler.GetAthleteAnalytics).Methods("POST")

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")

	router := middleware.LoggerMiddleware(mux)
	router = middleware.AuthMiddleware(enforcer, router)

	return router
}

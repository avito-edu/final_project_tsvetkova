package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"swim_service/internal/dto"
	"swim_service/internal/handler"
	"swim_service/internal/repository/athlete"
	"swim_service/internal/repository/competition"
	"swim_service/internal/repository/organizer"
	"swim_service/internal/repository/user"
	"swim_service/internal/service"
	"swim_service/pkg/jwtutil"

	"github.com/casbin/casbin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://postgres:password@localhost:5435/swim_service?sslmode=disable"
	}

	serverPORT := os.Getenv("SERVER_PORT")
	if serverPORT == "" {
		serverPORT = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	authEnforcer, err := casbin.NewEnforcerSafe("./config/auth_model.conf", "./config/policy.csv")
	if err != nil {
		log.Fatalf("failed to initialize authorization enforcer: %s", err.Error())
	}

	var db *sql.DB

	for i := range 5 {
		db, err = sql.Open("pgx", postgresURL)
		if err != nil {
			log.Printf("Failed to connect to database (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Database ping failed (attempt %d): %v", i+1, err)
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}
	defer db.Close()

	if err := runMigrationsFromFile(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	athleteRepo := athlete.NewAthleteRepository(db)
	organizerRepo := organizer.NewOrganizerRepository(db)
	competitionRepo := competition.NewCompetitionRepository(db)
	userRepo := user.NewUserRepository(db)

	tokenService := jwtutil.NewTokenService([]byte(jwtSecret))

	athleteService := service.NewAthleteService(athleteRepo)
	organizerService := service.NewOrganizerService(organizerRepo)
	competitionService := service.NewCompetitionService(competitionRepo)
	analyticsService := service.NewAnalyticsService(competitionRepo)
	userService := service.NewUserService(userRepo, tokenService)

	athleteHandler := handler.NewAthleteHandler(athleteService)
	organizerHandler := handler.NewOrganizerHandler(organizerService)
	competitionHandler := handler.NewCompetitionHandler(competitionService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	userHandler := handler.NewUserHandler(userService)

	router := handler.NewRouter(athleteHandler, competitionHandler, organizerHandler, analyticsHandler, userHandler, authEnforcer)

	admin := dto.RegisterSpecialRequest{
		Login:    "Admin",
		Password: "AdminPassword123",
		Role:     "admin",
	}

	_, err = userService.RegisterSpecial(context.Background(), admin)
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Addr:    ":" + serverPORT,
		Handler: router,
	}
	go func() {
		log.Printf("Server starting on port %s", serverPORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")

}

func runMigrationsFromFile(db *sql.DB) error {
	migrationsDir := filepath.Join(".", "migrations")
	
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = filepath.Join("..", "migrations")
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return fmt.Errorf("not found :(")
		}
	}

	goose.SetBaseFS(nil)
	
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
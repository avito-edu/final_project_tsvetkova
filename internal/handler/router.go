package handler

import (
	"swim_service/internal/middleware"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	athleteHandler *AthleteHandler,
	competitionHandler *CompetitionHandler,
	organizerHandler *OrganizerHandler,
	analyticsHandler *AnalyticsHandler,
	userHandler *UserHandler,
	enforcer *casbin.Enforcer,
) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.AuthMiddleware(enforcer))

	athleteGroup := router.Group("/athletes")
	{
		athleteGroup.POST("/create", athleteHandler.CreateAthlete)
		athleteGroup.GET("/get/:id", athleteHandler.GetAthlete)
		athleteGroup.PATCH("/update/:id", athleteHandler.UpdateStatus)
	}

	competitionGroup := router.Group("/competitions")
	{
		competitionGroup.POST("/create", competitionHandler.CreateCompetition)
		competitionGroup.GET("/get/:id", competitionHandler.GetCompetition)
		competitionGroup.GET("/get-all", competitionHandler.GetAllCompetitions)
		competitionGroup.POST("/add/:id", competitionHandler.AddResult)
	}

	organizerGroup := router.Group("/organizers")
	{
		organizerGroup.POST("/create", organizerHandler.CreateOrganizer)
		organizerGroup.GET("/get/:id", organizerHandler.GetOrganizer)
	}

	userGroup := router.Group("")
	{
		userGroup.POST("/register", userHandler.RegisterUser)
		userGroup.POST("/register-special", userHandler.RegisterSpecial)
		userGroup.POST("/login", userHandler.Login)
	}

	analyticsGroup := router.Group("/analytics")
	{
		analyticsGroup.POST("/get/:id", analyticsHandler.GetAthleteAnalytics)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
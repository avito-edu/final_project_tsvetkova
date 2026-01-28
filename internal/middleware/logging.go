package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initLogrus() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile("logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.SetOutput(os.Stdout)
		logger.WithError(err).Warn("Failed to open log file, using stdout")
	} else {
		logger.SetOutput(file)
	}

	return logger
}

func LoggerMiddleware() gin.HandlerFunc {
	logger := initLogrus()

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		logger.WithFields(logrus.Fields{
			"time":     end.Format("2006-01-02 15:04:05"),
			"method":   c.Request.Method,
			"path":     path,
			"query":    query,
			"status":   c.Writer.Status(),
			"duration": latency.Milliseconds(),
			"ip":       c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
		}).Info("request")
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

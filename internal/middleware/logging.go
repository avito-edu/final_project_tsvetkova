package middleware

import (
	"net/http"
	"os"
	"time"

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

func LoggerMiddleware(next http.Handler) http.Handler {
	logger := initLogrus()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseRecorder{
			ResponseWriter: w,
			status:         200,
		}

		next.ServeHTTP(rw, r)

		logger.WithFields(logrus.Fields{
			"time":     time.Now().Format("2006-01-02 15:04:05"),
			"method":   r.Method,
			"path":     r.URL.Path,
			"query":    r.URL.RawQuery,
			"status":   rw.status,
			"duration": time.Since(start).Milliseconds(),
			"ip":       r.RemoteAddr,
		}).Info("request")
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

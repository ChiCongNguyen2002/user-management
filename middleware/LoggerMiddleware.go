package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs details of each request including method, path, status, and latency.
func LoggerMiddleware() gin.HandlerFunc {
	logFile, err := os.OpenFile("middleware/logger.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(logFile)

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		log.Printf("[%s] %s %s - Status: %d - Latency: %v\n",
			start.Format("2006/01/02 - 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
}

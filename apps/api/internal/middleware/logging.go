package middleware

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LoggingMiddleware provides request/response logging
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("X-Request-ID", requestID)
		c.Locals("request_id", requestID)

		// Capture start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request details
		logRequest(c, requestID, start, duration, err)

		return err
	}
}

// logRequest logs request details
func logRequest(c *fiber.Ctx, requestID string, start time.Time, duration time.Duration, err error) {
	method := c.Method()
	path := c.Route().Path
	status := c.Response().StatusCode()
	clientIP := c.IP()
	userAgent := c.Get("User-Agent")
	userID := c.Locals("user_id")
	
	// Build log entry
	logEntry := map[string]interface{}{
		"request_id":   requestID,
		"timestamp":     start.Format(time.RFC3339),
		"duration_ms":   duration.Milliseconds(),
		"method":        method,
		"path":          path,
		"status":        status,
		"client_ip":     clientIP,
		"user_agent":    userAgent,
	}

	// Add user info if authenticated
	if userID != nil {
		logEntry["user_id"] = userID
	}

	// Add error if present
	if err != nil {
		logEntry["error"] = err.Error()
	}

	// Add query parameters if any
	if len(c.Queries()) > 0 {
		logEntry["query_params"] = c.Queries()
	}

	// Add response size
	if c.Response().Header.ContentLength() > 0 {
		logEntry["response_size_bytes"] = c.Response().Header.ContentLength()
	}

	// For now, just print to console (in production, use structured logging)
	logJSON(logEntry)
}

// logJSON outputs JSON log
func logJSON(data interface{}) {
	// In production, use proper logging library like logrus or zap
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Fallback to simple logging
		println("Failed to marshal log data:", err.Error())
		return
	}
	println(string(jsonData)) // Using println for simplicity - replace with proper logger
}

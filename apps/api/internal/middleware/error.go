package middleware

import (
	"api/internal/errors"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// ErrorMiddleware handles application errors consistently
func ErrorMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Capture any errors that occur during request processing
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\nStack: %s", err, debug.Stack())
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
					"code":  fiber.StatusInternalServerError,
				})
			}
		}()

		// Continue to next middleware
		return c.Next()
	}
}

// ErrorHandler is a global error handler that catches all errors
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Check if it's an AppError
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(fiber.Map{
			"error": appErr.Message,
			"code":  appErr.Code,
			"details": appErr.Details,
		})
	}

	// Handle Fiber errors
	if fiberErr, ok := err.(*fiber.Error); ok {
		// Don't log client errors (4xx)
		if fiberErr.Code < 500 {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"error": fiberErr.Message,
				"code":  fiberErr.Code,
			})
		}

		// Log server errors (5xx)
		log.Printf("Fiber Error: %s\nStack: %s", fiberErr.Message, debug.Stack())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// Handle other unexpected errors
	if err != nil {
		// Don't log context cancellation errors
		if err.Error() == "context canceled" {
			return c.SendStatus(fiber.StatusRequestTimeout)
		}

		// Log unexpected errors
		log.Printf("Unexpected Error: %s\nStack: %s", err.Error(), debug.Stack())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
			"code":  fiber.StatusInternalServerError,
		})
	}

	// No error
	return nil
}

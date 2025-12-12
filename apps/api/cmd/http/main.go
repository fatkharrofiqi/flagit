package main

import (
	database "api/internal/config/db"
	"api/internal/config/env"
	"api/internal/controller"
	"api/internal/middleware"
	"api/internal/repository"
	"api/internal/route"
	"api/internal/service"
	"api/internal/validation"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := env.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.InitializeDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	projectRepo := repository.NewProjectRepository(db)
	envRepo := repository.NewEnvironmentRepository(db)
	flagRepo := repository.NewFlagRepository(db)
	flagValueRepo := repository.NewFlagValueRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize SSE controller
	sseController := controller.NewSSEController()

	// Initialize validator once
	validator := validation.NewValidator()

	// Get JWT secret from config (in production, use environment variable)
	jwtSecret := "your-super-secret-jwt-key-change-in-production" // TODO: Get from config

	// Initialize services with SSE controller
	projectService := service.NewProjectService(projectRepo, sseController)
	envService := service.NewEnvironmentService(envRepo, sseController)
	flagService := service.NewFlagService(flagRepo, flagValueRepo, sseController)
	authService := service.NewAuthService(userRepo, jwtSecret)

	// Initialize controllers
	projectController := controller.NewProjectController(projectService, validator)
	envController := controller.NewEnvironmentController(envService)
	flagController := controller.NewFlagController(flagService)
	authController := controller.NewAuthController(authService, validator)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.LoggingMiddleware())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	
	// Custom error middleware (must be last)
	app.Use(middleware.ErrorMiddleware())

	// Setup routes
	router := route.NewRouter(app, projectController, envController, flagController, sseController, authController, cfg)
	router.SetupRoutes()

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": c.Context().Time(),
		})
	})

	// Start server
	port := cfg.Server.Port
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

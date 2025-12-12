package route

import (
	"api/internal/controller"
	"api/internal/config/env"
	"api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app                 *fiber.App
	projectController   *controller.ProjectController
	envController       *controller.EnvironmentController
	flagController      *controller.FlagController
	sseController       *controller.SSEController
	authController      *controller.AuthController
}

func NewRouter(
	app *fiber.App,
	projectController *controller.ProjectController,
	envController *controller.EnvironmentController,
	flagController *controller.FlagController,
	sseController *controller.SSEController,
	authController *controller.AuthController,
	cfg *env.Config,
) *Router {
	router := &Router{
		app:                 app,
		projectController:   projectController,
		envController:       envController,
		flagController:      flagController,
		sseController:       sseController,
		authController:      authController,
	}
	
	// Store JWT secret in router struct
	// In a real app, you'd want to handle this more securely
	return router
}

func (r *Router) SetupRoutes() {
	api := r.app.Group("/api")

	// Authentication routes (public)
	api.Post("/auth/register", r.authController.Register)
	api.Post("/auth/login", r.authController.Login)
	
	// Profile route (authenticated)
	api.Get("/auth/profile", r.authController.Profile)

	// SSE endpoint for real-time updates
	api.Get("/events", r.sseController.RegisterClient)

	// Projects (secured with authentication)
	projects := api.Group("/projects")
	projects.Use(middleware.AuthMiddleware("jwt-secret-placeholder")) // TODO: Get from config
	projects.Get("/", middleware.RequirePermission(middleware.ProjectRead), r.projectController.GetProjects)
	projects.Post("/", middleware.RequirePermission(middleware.ProjectCreate), r.projectController.CreateProject)
	projects.Get("/:id", middleware.RequirePermission(middleware.ProjectRead), r.projectController.GetProject)
	projects.Put("/:id", middleware.RequirePermission(middleware.ProjectUpdate), r.projectController.UpdateProject)
	projects.Delete("/:id", middleware.RequirePermission(middleware.ProjectDelete), r.projectController.DeleteProject)

	// Environments (secured)
	environments := api.Group("/environments")
	environments.Use(middleware.AuthMiddleware("jwt-secret-placeholder")) // TODO: Get from config
	environments.Get("/", middleware.RequirePermission(middleware.EnvironmentRead), r.envController.GetEnvironments)
	environments.Post("/", middleware.RequirePermission(middleware.EnvironmentCreate), r.envController.CreateEnvironment)
	environments.Get("/:id", middleware.RequirePermission(middleware.EnvironmentRead), r.envController.GetEnvironment)
	environments.Put("/:id", middleware.RequirePermission(middleware.EnvironmentUpdate), r.envController.UpdateEnvironment)
	environments.Delete("/:id", middleware.RequirePermission(middleware.EnvironmentDelete), r.envController.DeleteEnvironment)

	// Project environments
	projects.Get("/:projectId/environments", middleware.RequirePermission(middleware.EnvironmentRead), r.envController.GetProjectEnvironments)

	// Flags (secured)
	flags := api.Group("/flags")
	flags.Use(middleware.AuthMiddleware("jwt-secret-placeholder")) // TODO: Get from config
	flags.Get("/", middleware.RequirePermission(middleware.FlagRead), r.flagController.GetProjectFlags) // Using GetProjectFlags as fallback
	flags.Post("/", middleware.RequirePermission(middleware.FlagCreate), r.flagController.CreateFlag)
	flags.Get("/:id", middleware.RequirePermission(middleware.FlagRead), r.flagController.GetFlag)
	flags.Put("/:id", middleware.RequirePermission(middleware.FlagUpdate), r.flagController.UpdateFlag)
	flags.Delete("/:id", middleware.RequirePermission(middleware.FlagDelete), r.flagController.DeleteFlag)

	// Project flags
	projects.Get("/:projectId/flags", middleware.RequirePermission(middleware.FlagRead), r.flagController.GetProjectFlags)
	
	// Flag values
	flags.Get("/:flagId/values", middleware.RequirePermission(middleware.FlagRead), r.flagController.GetFlagValues)
	flags.Post("/values", middleware.RequirePermission(middleware.FlagCreate), r.flagController.CreateOrUpdateFlagValue)
	flags.Put("/values/:id", middleware.RequirePermission(middleware.FlagUpdate), r.flagController.UpdateFlagValue)
	flags.Delete("/values/:id", middleware.RequirePermission(middleware.FlagDelete), r.flagController.DeleteFlagValue)
	
	// Environment flags
	environments.Get("/:envId/flags", middleware.RequirePermission(middleware.FlagRead), r.flagController.GetEnvironmentFlags)
}

package controller

import (
	"api/internal/dto"
	"api/internal/service"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EnvironmentController struct {
	service service.EnvironmentService
}

func NewEnvironmentController(service service.EnvironmentService) *EnvironmentController {
	return &EnvironmentController{service: service}
}

func (c *EnvironmentController) CreateEnvironment(ctx *fiber.Ctx) error {
	var req dto.CreateEnvironmentRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	env, err := c.service.CreateEnvironment(context.Background(), &req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create environment",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(env)
}

func (c *EnvironmentController) GetEnvironments(ctx *fiber.Ctx) error {
	environments, err := c.service.GetAllEnvironments(context.Background())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch environments",
		})
	}

	return ctx.JSON(environments)
}

func (c *EnvironmentController) GetProjectEnvironments(ctx *fiber.Ctx) error {
	projectIDStr := ctx.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	environments, err := c.service.GetProjectEnvironments(context.Background(), projectID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch project environments",
		})
	}

	return ctx.JSON(environments)
}

func (c *EnvironmentController) GetEnvironment(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid environment ID",
		})
	}

	env, err := c.service.GetEnvironment(context.Background(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch environment",
		})
	}

	return ctx.JSON(env)
}

func (c *EnvironmentController) UpdateEnvironment(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid environment ID",
		})
	}

	var req dto.UpdateEnvironmentRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	env, err := c.service.UpdateEnvironment(context.Background(), id, &req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update environment",
		})
	}

	return ctx.JSON(env)
}

func (c *EnvironmentController) DeleteEnvironment(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid environment ID",
		})
	}

	err = c.service.DeleteEnvironment(context.Background(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete environment",
		})
	}

	return ctx.Status(http.StatusNoContent).Send(nil)
}

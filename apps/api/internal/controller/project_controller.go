package controller

import (
	"api/internal/dto"
	"api/internal/service"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProjectController struct {
	service service.ProjectService
}

func NewProjectController(service service.ProjectService) *ProjectController {
	return &ProjectController{service: service}
}

func (c *ProjectController) CreateProject(ctx *fiber.Ctx) error {
	var req dto.CreateProjectRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	project, err := c.service.CreateProject(context.Background(), &req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create project",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(project)
}

func (c *ProjectController) GetProjects(ctx *fiber.Ctx) error {
	projects, err := c.service.GetAllProjects(context.Background())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch projects",
		})
	}

	return ctx.JSON(projects)
}

func (c *ProjectController) GetProject(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	project, err := c.service.GetProject(context.Background(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch project",
		})
	}

	return ctx.JSON(project)
}

func (c *ProjectController) UpdateProject(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	var req dto.UpdateProjectRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	project, err := c.service.UpdateProject(context.Background(), id, &req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update project",
		})
	}

	return ctx.JSON(project)
}

func (c *ProjectController) DeleteProject(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	err = c.service.DeleteProject(context.Background(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete project",
		})
	}

	return ctx.Status(http.StatusNoContent).Send(nil)
}

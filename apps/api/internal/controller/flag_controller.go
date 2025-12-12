package controller

import (
	"api/internal/dto"
	"api/internal/service"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FlagController struct {
	service service.FlagService
}

func NewFlagController(service service.FlagService) *FlagController {
	return &FlagController{service: service}
}

func (c *FlagController) CreateFlag(ctx *fiber.Ctx) error {
	var req dto.CreateFlagRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	flag, err := c.service.CreateFlag(context.Background(), &req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create flag",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(flag)
}

func (c *FlagController) GetProjectFlags(ctx *fiber.Ctx) error {
	projectIDStr := ctx.Params("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	includeValues := ctx.Query("includeValues") == "true"

	flags, err := c.service.GetProjectFlags(context.Background(), projectID, includeValues)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch project flags",
		})
	}

	return ctx.JSON(flags)
}

func (c *FlagController) GetFlag(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid flag ID",
		})
	}

	flag, err := c.service.GetFlag(context.Background(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch flag",
		})
	}

	return ctx.JSON(flag)
}

func (c *FlagController) UpdateFlag(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid flag ID",
		})
	}

	var req dto.UpdateFlagRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	flag, err := c.service.UpdateFlag(context.Background(), id, &req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update flag",
		})
	}

	return ctx.JSON(flag)
}

func (c *FlagController) DeleteFlag(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid flag ID",
		})
	}

	err = c.service.DeleteFlag(context.Background(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete flag",
		})
	}

	return ctx.Status(http.StatusNoContent).Send(nil)
}

// Flag value endpoints
func (c *FlagController) GetFlagValues(ctx *fiber.Ctx) error {
	flagIDStr := ctx.Params("flagId")
	flagID, err := uuid.Parse(flagIDStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid flag ID",
		})
	}

	values, err := c.service.GetFlagValues(context.Background(), flagID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch flag values",
		})
	}

	return ctx.JSON(values)
}

func (c *FlagController) GetEnvironmentFlags(ctx *fiber.Ctx) error {
	envIDStr := ctx.Params("envId")
	envID, err := uuid.Parse(envIDStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid environment ID",
		})
	}

	values, err := c.service.GetEnvironmentFlags(context.Background(), envID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch environment flags",
		})
	}

	return ctx.JSON(values)
}

func (c *FlagController) CreateOrUpdateFlagValue(ctx *fiber.Ctx) error {
	var req dto.CreateFlagValueRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	flagValue, err := c.service.CreateOrUpdateFlagValue(context.Background(), &req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create/update flag value",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(flagValue)
}

func (c *FlagController) UpdateFlagValue(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid flag value ID",
		})
	}

	var req dto.UpdateFlagValueRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	flagValue, err := c.service.UpdateFlagValue(context.Background(), id, &req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update flag value",
		})
	}

	return ctx.JSON(flagValue)
}

func (c *FlagController) DeleteFlagValue(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid flag value ID",
		})
	}

	err = c.service.DeleteFlagValue(context.Background(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete flag value",
		})
	}

	return ctx.Status(http.StatusNoContent).Send(nil)
}

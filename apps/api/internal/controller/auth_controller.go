package controller

import (
	"api/internal/errors"
	"api/internal/model"
	"api/internal/service"
	"api/internal/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest) // Error handled by middleware
	}

	// Initialize validator
	validator := validation.NewValidator()

	// Validate request fields
	if err := validator.ValidateUsername(req.Username); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := validator.ValidatePassword(req.Password); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := validator.ValidateName(req.FirstName); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := validator.ValidateName(req.LastName); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	response, err := c.authService.Register(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors.NewAppError(fiber.StatusBadRequest, err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req model.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors.ErrInvalidRequestBody)
	}

	// Initialize validator
	validator := validation.NewValidator()

	// Validate request fields
	if err := validator.ValidateUsername(req.Username); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	if err := validator.ValidatePassword(req.Password); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	response, err := c.authService.Login(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errors.ErrInvalidCredentials)
	}

	return ctx.JSON(response)
}

func (c *AuthController) Profile(ctx *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr := ctx.Locals("user_id")
	if userIDStr == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errors.ErrAuthenticationRequired)
	}

	// Convert string to UUID (using google/uuid)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors.ErrInvalidField)
	}

	user, err := c.authService.GetUserByID(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.ErrInternalServer)
	}

	if user == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(errors.ErrUserNotFound)
	}

	return ctx.JSON(user.ToResponse())
}

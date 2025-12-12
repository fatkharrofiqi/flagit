package middleware

import (
	"api/internal/errors"

	"github.com/gofiber/fiber/v2"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleManager  Role = "manager"
	RoleDeveloper Role = "developer"
	RoleViewer   Role = "viewer"
)

// Permission represents specific actions users can perform
type Permission string

const (
	// Project permissions
	ProjectCreate Permission = "project:create"
	ProjectRead   Permission = "project:read"
	ProjectUpdate Permission = "project:update"
	ProjectDelete Permission = "project:delete"

	// Environment permissions
	EnvironmentCreate Permission = "environment:create"
	EnvironmentRead   Permission = "environment:read"
	EnvironmentUpdate Permission = "environment:update"
	EnvironmentDelete Permission = "environment:delete"

	// Flag permissions
	FlagCreate Permission = "flag:create"
	FlagRead   Permission = "flag:read"
	FlagUpdate Permission = "flag:update"
	FlagDelete Permission = "flag:delete"
)

// RolePermissions maps roles to their allowed permissions
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		// Admin has all permissions
		ProjectCreate, ProjectRead, ProjectUpdate, ProjectDelete,
		EnvironmentCreate, EnvironmentRead, EnvironmentUpdate, EnvironmentDelete,
		FlagCreate, FlagRead, FlagUpdate, FlagDelete,
	},
	RoleManager: {
		// Manager can manage everything within their projects
		ProjectCreate, ProjectRead, ProjectUpdate,
		EnvironmentCreate, EnvironmentRead, EnvironmentUpdate, EnvironmentDelete,
		FlagCreate, FlagRead, FlagUpdate, FlagDelete,
	},
	RoleDeveloper: {
		// Developer can read and create/update flags
		ProjectRead,
		EnvironmentRead, EnvironmentCreate,
		FlagCreate, FlagRead, FlagUpdate,
	},
	RoleViewer: {
		// Viewer can only read
		ProjectRead,
		EnvironmentRead,
		FlagRead,
	},
}

// HasPermission checks if a user role has a specific permission
func HasPermission(role Role, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}

	return false
}

// RequirePermission creates middleware that checks if user has required permission
func RequirePermission(permission Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by auth middleware)
		roleStr, ok := c.Locals("role").(string)
		if !ok || roleStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.ErrAuthenticationRequired)
		}

		role := Role(roleStr)
		if !HasPermission(role, permission) {
			return c.Status(fiber.StatusForbidden).JSON(errors.ErrInsufficientPermissions)
		}

		return c.Next()
	}
}

// RequireRole creates middleware that checks if user has required role
func RequireRole(requiredRole Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by auth middleware)
		roleStr, ok := c.Locals("role").(string)
		if !ok || roleStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.ErrAuthenticationRequired)
		}

		role := Role(roleStr)
		if role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(errors.ErrInsufficientRole)
		}

		return c.Next()
	}
}

// RequireAdmin creates middleware that requires admin role
func RequireAdmin() fiber.Handler {
	return RequireRole(RoleAdmin)
}

// RequireOwnerOrAdmin creates middleware that checks if user is resource owner or admin
func RequireOwnerOrAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user info from context
		userID, userOk := c.Locals("user_id").(string)
		roleStr, roleOk := c.Locals("role").(string)
		if !userOk || !roleOk || userID == "" || roleStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.ErrAuthenticationRequired)
		}

		role := Role(roleStr)
		// Admin can access everything
		if role == RoleAdmin {
			return c.Next()
		}

		// TODO: Implement ownership checking based on resource
		// For now, we'll check path parameter against user ID
		resourceID := c.Params("id")
		if userID == resourceID {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(errors.ErrAccessDenied)
	}
}

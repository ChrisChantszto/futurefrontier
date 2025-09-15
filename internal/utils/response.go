package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/models"
)

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(message, data))
}

// RespondWithError sends a standardized error response
func RespondWithError(c *fiber.Ctx, message string, errorCode int, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(models.NewErrorResponse(message, errorCode, data))
}

// RespondWithServerError sends a standardized server error response
func RespondWithServerError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
		"Server error occurred",
		models.ErrAppServer,
		map[string]interface{}{
			"error": err.Error(),
		},
	))
}

// Common response messages
const (
	MsgDataFound            = "Data found"
	MsgNoDataFound          = "No data found"
	MsgNotFound             = "Not found"
	MsgUpdatedSuccessfully  = "Updated successfully"
	MsgCreatedSuccessfully  = "Created successfully"
	MsgDeletedSuccessfully  = "Deleted successfully"
	MsgLoginSuccess         = "User logged in successfully"
	MsgLogoutSuccess        = "User logged out successfully"
	MsgInvalidCredentials   = "Invalid email or password"
	MsgUnauthorized         = "You are not authorized to perform this action"
	MsgForbidden            = "Access forbidden"
)

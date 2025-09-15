package handlers

import "github.com/gofiber/fiber/v2"

// JSONSuccess returns a standardized success response
func JSONSuccess(c *fiber.Ctx, message string, data any) error {
	if message == "" {
		message = "success"
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// JSONSuccessWithExtra returns a standardized success response and allows extra top-level fields
func JSONSuccessWithExtra(c *fiber.Ctx, message string, data any, extra fiber.Map) error {
	if message == "" {
		message = "success"
	}
	resp := fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	}
	for k, v := range extra {
		resp[k] = v
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// JSONError returns a standardized error response with provided HTTP status
// Note: We keep proper HTTP status codes to avoid breaking existing frontend error handling.
func JSONError(c *fiber.Ctx, status int, message string, errorCode int, data any) error {
	if status <= 0 {
		status = fiber.StatusBadRequest
	}
	if message == "" {
		message = "error"
	}
	resp := fiber.Map{
		"success":    false,
		"message":    message,
		"data":       data,
	}
	if errorCode != 0 {
		resp["error_code"] = errorCode
	}
	return c.Status(status).JSON(resp)
}

// JSONErrorWithExtra returns an error response and allows extra top-level fields
func JSONErrorWithExtra(c *fiber.Ctx, status int, message string, errorCode int, data any, extra fiber.Map) error {
	if status <= 0 {
		status = fiber.StatusBadRequest
	}
	if message == "" {
		message = "error"
	}
	resp := fiber.Map{
		"success": false,
		"message": message,
		"data":    data,
	}
	if errorCode != 0 {
		resp["error_code"] = errorCode
	}
	for k, v := range extra {
		resp[k] = v
	}
	return c.Status(status).JSON(resp)
}

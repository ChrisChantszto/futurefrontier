package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onetakesolutions/onetake-corpsite-backend/internal/config"
)

// RequireAuth validates the "access" JWT cookie and sets c.Locals("userEmail", email)
func RequireAuth(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tok := c.Cookies("access")
		if tok == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing access token")
		}
		parsed, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !parsed.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid access token")
		}
		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid claims")
		}
		email, ok := claims["email"].(string)
		if !ok || email == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user")
		}
		c.Locals("userEmail", email)
		return c.Next()
	}
}

package middleware

import (
	"strings"

	"file_project/utils"

	"github.com/gofiber/fiber/v2"
)

// JWTProtected protects routes using Authorization: Bearer <token>
func JWTProtected(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header"})
	}
	claims, err := utils.ParseJWT(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
	}
	// store user info in locals for handlers
	c.Locals("user_id", claims.UserID)
	c.Locals("email", claims.Email)
	return c.Next()
}

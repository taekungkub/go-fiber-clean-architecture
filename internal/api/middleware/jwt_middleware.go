package middleware

import (
	"go-fiber-api/pkg/core"
	"go-fiber-api/pkg/dto"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware validates JWT tokens and protects routes
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return dto.SendError(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return dto.SendError(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := parts[1]

		// Validate token
		claims, err := core.ValidateToken(tokenString)
		if err != nil {
			return dto.SendError(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

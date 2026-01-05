package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CurrentAuthUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		t, ok := c.Locals("jwt").(*jwt.Token)
		if !ok || t == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		id, _ := claims["id"].(string)
		email, _ := claims["email"].(string)

		c.Locals("userID", id)
		c.Locals("userEmail", email)
		return c.Next()
	}
}

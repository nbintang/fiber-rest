package middleware

import (
	"rest-fiber/internal/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const AccessAuth string = "access-auth"

func CurrentAuthUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals(AccessAuth).(*jwt.Token)
		if !ok || token == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		id, _ := claims["id"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(enums.EUserRoleType)
		c.Locals("userID", id)
		c.Locals("userEmail", email)
		c.Locals("userRole", role)
		return c.Next()
	}
}

package middleware

import (
	"rest-fiber/internal/enums" 

	"github.com/gofiber/fiber/v2"
)


type AuthClaims struct {
	ID    string
	Email string
	Role  enums.EUserRoleType
	JTI   string
}


func GetCurrentUser(c *fiber.Ctx) (*AuthClaims, error) {
	v := c.Locals(enums.CurrentUserKey)
	user, ok := v.(*AuthClaims);
	if !ok ||user == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	};
	return user, nil
}
package middleware

import (
	"rest-fiber/utils/enums"

	"github.com/gofiber/fiber/v2"
)

func AuthRoleAccess(roles ...enums.EUserRoleType) fiber.Handler {
	roleSet := make(map[enums.EUserRoleType]struct{}, len(roles))
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		user, err := GetCurrentUser(c)
		if err != nil {
			return err
		}
		if _, ok := roleSet[user.Role]; !ok {
			return fiber.NewError(fiber.StatusForbidden, "Forbidden")
		}
		return c.Next()
	}
}

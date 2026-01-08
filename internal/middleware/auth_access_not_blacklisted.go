package middleware

import (
	"rest-fiber/internal/infra"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func AccessNotBlacklisted(redisService infra.RedisService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals(AccessAuth).(*jwt.Token)
		if !ok || token == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		jti, _ := claims["jti"].(string)
		if jti == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		key := "blacklist_access:" + jti
		_, err := redisService.Get(c.UserContext(), key)
		if err == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "token revoked")
		}
		if err != redis.Nil {
			return fiber.NewError(fiber.StatusInternalServerError, "redis error")
		}
		return c.Next()
	}
}

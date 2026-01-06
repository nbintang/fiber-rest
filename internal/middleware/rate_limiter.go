package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage"
)

type RateLimiterParams struct {
	MaxLimit int
	Storage  storage.Storage
}

func RateLimiter(params RateLimiterParams) fiber.Handler {
	config := limiter.Config{
		Max:        params.MaxLimit,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Too many requests")
		},
	}
	if params.Storage != nil {
		config.Storage = params.Storage
	}
	return limiter.New(config)
}

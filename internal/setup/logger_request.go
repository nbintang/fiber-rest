package setup

import (
	"rest-fiber/internal/infra" 

	"github.com/gofiber/fiber/v2" 
)

func LoggerRequest(l *infra.AppLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		l.Info("http_request")
		return err
	}
}

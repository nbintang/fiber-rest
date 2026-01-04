package setup

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func DefaultErrorHandler(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError
	msg := "Internal Server Error"
	if ve, ok := err.(validator.ValidationErrors); ok {
		out := make([]fiber.Map, 0, len(ve))
		for _, fe := range ve {
			out = append(out, fiber.Map{
				"field": fe.Field(),
				"tag":   fe.Tag(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation error",
			"errors":  out,
		})
	}

	if e, ok := err.(*fiber.Error); ok {
		statusCode = e.Code
		msg = e.Message
	}
	return c.Status(statusCode).JSON(fiber.Map{
		"error":     msg,
		"status":    statusCode,
		"timestamp": time.Now().Unix(),
	})
}

package setup

import (
	"rest-fiber/pkg/httpx"

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
		return c.Status(fiber.StatusBadRequest).JSON(httpx.NewHttpResponse(statusCode, msg, out))
	}

	if e, ok := err.(*fiber.Error); ok {
		statusCode = e.Code
		msg = e.Message
	}
	return c.Status(statusCode).JSON(httpx.NewHttpResponse[any](statusCode, msg, nil))
}

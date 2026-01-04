package contract

import "github.com/gofiber/fiber/v2"

type Route interface {
	RegisterRoute(api fiber.Router)
}
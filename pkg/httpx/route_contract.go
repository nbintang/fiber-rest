package httpx

import "github.com/gofiber/fiber/v2"

type Route interface {
	RegisterRoute(api fiber.Router)
}

type ProtectedRoute interface {
	RegisterProtectedRoute(api fiber.Router)
}
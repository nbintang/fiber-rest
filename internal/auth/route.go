package auth

import (
	"rest-fiber/internal/contract"

	"github.com/gofiber/fiber/v2"
)

type authRouteImpl struct {
	h AuthHandler
}

func NewAuthRoute(h AuthHandler) contract.Route {
	return &authRouteImpl{h}
}
func (r *authRouteImpl) RegisterRoute(api fiber.Router) {
	auth := api.Group("/auth")
	auth.Post("/register", r.h.Register)
}

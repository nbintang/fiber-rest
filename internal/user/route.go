package user

import (
	"rest-fiber/internal/contract"

	"github.com/gofiber/fiber/v2"
)

type userRouteImpl struct {
	h UserHandler
}

func NewUserRoute(h UserHandler) contract.Route {
	return &userRouteImpl{h}
}
func (r *userRouteImpl) RegisterRoute(api fiber.Router) {
	users := api.Group("/users")
	users.Get("/", r.h.GetAllUsers)
	users.Get("/:id", r.h.GetUserByID)
} 



package user

import (
	"rest-fiber/internal/middleware"
	"rest-fiber/pkg/httpx"
	"rest-fiber/utils/enums"

	"github.com/gofiber/fiber/v2"
)


type UserRouteParams struct{
	httpx.RouteParams
	UserHandler UserHandler
}
type userRouteImpl struct {
	userHandler UserHandler
}

func NewUserRoute(params UserRouteParams) httpx.ProtectedRoute {
	return &userRouteImpl{userHandler: params.UserHandler}
}
func (r *userRouteImpl) RegisterProtectedRoute(api fiber.Router) {
	users := api.Group("/users")
	users.Get("/me", r.userHandler.GetCurrentUserProfile)
	users.Patch("/me", r.userHandler.UpdateCurrentUser)
	users.Get("/", middleware.AuthRoleAccess(enums.Admin), r.userHandler.GetAllUsers)
	users.Get("/:id", middleware.AuthRoleAccess(enums.Admin), r.userHandler.GetUserByID)
}

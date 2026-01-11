package category

import (
	"rest-fiber/internal/middleware"
	"rest-fiber/pkg/httpx"
	"rest-fiber/utils/enums"

	"github.com/gofiber/fiber/v2"
)

type CategoryRouteParams struct {
	httpx.RouteParams
	CategoryHandler CategoryHandler
}

type categoryRouteImpl struct {
	categoryHandler CategoryHandler
}

func NewCategoryRoutes(params CategoryRouteParams) httpx.ProtectedRoute {
	return &categoryRouteImpl{categoryHandler: params.CategoryHandler}
}

func (r *categoryRouteImpl) RegisterProtectedRoute(route fiber.Router) {
	categories := route.Group("/categories")
	categories.Get("/", r.categoryHandler.GetAllCategories)
	categories.Get("/:id", r.categoryHandler.GetCategoryByID)
	categories.Post("/", middleware.AuthAllowRoleAccess(enums.Admin), r.categoryHandler.CreateCategory)
	categories.Patch("/:id", middleware.AuthAllowRoleAccess(enums.Admin), r.categoryHandler.UpdateCategoryByID)
	categories.Delete("/:id", middleware.AuthAllowRoleAccess(enums.Admin), r.categoryHandler.DeleteCategoryByID)
}

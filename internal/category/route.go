package category

import (
	"rest-fiber/internal/enums"

	"rest-fiber/internal/http/middleware"
	"rest-fiber/internal/http/router"

	"github.com/gofiber/fiber/v2"
)

type CategoryRouteParams struct {
	router.RouteParams
	CategoryHandler CategoryHandler
}

type categoryRouteImpl struct {
	categoryHandler CategoryHandler
}

func NewCategoryRoutes(params CategoryRouteParams) router.ProtectedRoute {
	return &categoryRouteImpl{categoryHandler: params.CategoryHandler}
}

func (r *categoryRouteImpl) RegisterProtectedRoute(route fiber.Router) {
	categories := route.Group("/categories")
	categories.Get("/", r.categoryHandler.GetAllCategories)
	categories.Get("/:id", r.categoryHandler.GetCategoryByID)
	categories.Post("/", middleware.AllowRoleAccess(enums.Admin), r.categoryHandler.CreateCategory)
	categories.Patch("/:id", middleware.AllowRoleAccess(enums.Admin), r.categoryHandler.UpdateCategoryByID)
	categories.Delete("/:id", middleware.AllowRoleAccess(enums.Admin), r.categoryHandler.DeleteCategoryByID)
}

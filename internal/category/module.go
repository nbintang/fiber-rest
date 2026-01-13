package category

import (
	"rest-fiber/internal/enums"
	"rest-fiber/internal/http/router" 

	"go.uber.org/fx"
)


var Module = fx.Module(
	"category",
	fx.Provide(
		NewCategoryRepository,
		NewCategoryService,
		NewCategoryHandler,
		router.ProvideRoute[CategoryRouteParams, router.ProtectedRoute](
			NewCategoryRoutes,
			enums.RouteProtected,
		),
	),
)
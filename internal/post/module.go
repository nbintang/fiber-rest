package post

import (
	"rest-fiber/internal/enums"
	"rest-fiber/internal/http/router" 
	"go.uber.org/fx"
)

var Module = fx.Module(
	"post",
	fx.Provide(
		NewPostRepository,
		NewPostService,
		NewPostHandler,
		router.ProvideRoute[PostRouteParams, router.ProtectedRoute](
			NewPostRoute,
			enums.RouteProtected,
		),
	),
)

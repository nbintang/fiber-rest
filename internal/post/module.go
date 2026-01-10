package post

import (
	"rest-fiber/internal/contract"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/setup"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"post",
	fx.Provide(
		NewPostRepository,
		NewPostService,
		NewPostHandler,
		setup.RouteProvider[PostRouteParams, contract.ProtectedRoute](
			NewPostRoute,
			enums.RouteProtected,
		),
	),
)

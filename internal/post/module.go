package post

import (
	"rest-fiber/pkg/httpx"
	"rest-fiber/utils/enums"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"post",
	fx.Provide(
		NewPostRepository,
		NewPostService,
		NewPostHandler,
		httpx.ProvideRoute[PostRouteParams, httpx.ProtectedRoute](
			NewPostRoute,
			enums.RouteProtected,
		),
	),
)

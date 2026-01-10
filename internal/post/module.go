package post

import (
	"rest-fiber/internal/contract"
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
		httpx.RouteProvider[PostRouteParams, contract.ProtectedRoute](
			NewPostRoute,
			enums.RouteProtected,
		),
	),
)

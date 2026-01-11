package user

import (
	"rest-fiber/pkg/httpx"
	"rest-fiber/utils/enums"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewUserRepository,
		NewUserService,
		NewUserHandler,
		httpx.ProvideRoute[UserRouteParams, httpx.ProtectedRoute](
			NewUserRoute,
			enums.RouteProtected,
		),
	),
)

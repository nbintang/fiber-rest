package user

import (
	"rest-fiber/internal/contract"
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
		httpx.RouteProvider[UserRouteParams, contract.ProtectedRoute](
			NewUserRoute,
			enums.RouteProtected,
		),
	),
)

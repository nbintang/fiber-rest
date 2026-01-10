package user

import (
	"rest-fiber/internal/contract"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/setup"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewUserRepository,
		NewUserService,
		NewUserHandler,
		setup.RouteProvider[UserRouteParams, contract.ProtectedRoute](
			NewUserRoute,
			enums.RouteProtected,
		),
	),
)

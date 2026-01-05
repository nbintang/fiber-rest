package user

import (
	"rest-fiber/internal/contract"
	"rest-fiber/internal/setup"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewUserRepository,
		NewUserService,
		NewUserHandler,
		setup.RouteProvider[UserHandler, contract.ProtectedRoute](
			NewUserRoute,
			setup.RouteProtected,
		),
	),
)

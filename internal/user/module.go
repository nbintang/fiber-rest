package user

import (
	"rest-fiber/internal/setup"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"User",
	fx.Provide(
		NewUserRepository,
		NewUserService,
		NewUserHandler,
		setup.RouteProvider[UserHandler](NewUserRoute),
	),
)

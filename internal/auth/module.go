package auth

import (
	"rest-fiber/internal/setup"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"Auth",
	fx.Provide(
		NewAuthService,
		NewAuthHandler,
		setup.RouteProvider[AuthHandler](NewAuthRoute),
	),
)

package auth

import (
	"rest-fiber/internal/contract"
	"rest-fiber/internal/setup"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		NewAuthService,
		NewAuthHandler,
		setup.RouteProvider[AuthHandler, contract.Route](
			NewAuthRoute,
			setup.RoutePublic,
		),
	),
)

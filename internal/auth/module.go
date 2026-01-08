package auth

import (
	"rest-fiber/internal/contract"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/setup"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		NewAuthService,
		NewAuthHandler,
		setup.RouteProvider[AuthRouteParams, contract.Route](
			NewAuthRoute,
			enums.RoutePublic,
		),
	),
)

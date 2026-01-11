package auth

import (
	"rest-fiber/pkg/httpx"
	"rest-fiber/utils/enums"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		NewAuthService,
		NewAuthHandler,
		httpx.ProvideRoute[AuthRouteParams, httpx.Route](
			NewAuthRoute,
			enums.RoutePublic,
		),
	),
)

package auth

import (
	"rest-fiber/internal/contract"
	"rest-fiber/pkg/httpx"
	"rest-fiber/utils/enums"

	"go.uber.org/fx"
) 

var Module = fx.Module(
	"auth",
	fx.Provide(
		NewAuthService,
		NewAuthHandler,
		httpx.RouteProvider[AuthRouteParams, contract.Route](
			NewAuthRoute,
			enums.RoutePublic,
		),
	),
)

package internal

import (
	"rest-fiber/internal/contract"
	"go.uber.org/fx"
)

type RoutesIn struct {
	fx.In
	App    *App
	Routes []contract.Route `group:"routes"`
}

func RegisterAllRoutes(in RoutesIn) {
	for _, r := range in.Routes {
		r.RegisterRoute(in.App.API)
	}
}

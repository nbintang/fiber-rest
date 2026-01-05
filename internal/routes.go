package internal

import (
	"rest-fiber/internal/contract"

	"go.uber.org/fx"
)

type RoutesIn struct {
	fx.In
	App       *App
	Routes    []contract.Route          `group:"public_routes"`
	Protected []contract.ProtectedRoute `group:"protected_routes"`
}

func RegisterAllRoutes(in RoutesIn) {
	for _, r := range in.Routes {
		r.RegisterRoute(in.App.PublicRoute)
	}
	for _, r := range in.Protected {
		r.RegisterProtectedRoute(in.App.ProtectedRoute)
	}
}

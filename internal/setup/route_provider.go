package setup

import "go.uber.org/fx"

type RouteConstructor[T any, R any] func(T) R

type AccessType string

const (
	RouteProtected AccessType = `group:"protected_routes"`
	RoutePublic    AccessType = `group:"public_routes"`
)

func RouteProvider[T any, R any](routeConstructor RouteConstructor[T, R], acc AccessType) any {
	return fx.Annotate(
		routeConstructor,
		fx.As(new(R)),
		fx.ResultTags(string(acc)),
	)
}

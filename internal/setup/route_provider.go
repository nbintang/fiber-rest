package setup

import (
	"rest-fiber/internal/enums"

	"go.uber.org/fx"
)

type RouteConstructor[T any, R any] func(T) R

func RouteProvider[T any, R any](routeConstructor RouteConstructor[T, R], acc enums.EAccessType) any {
	return fx.Annotate(
		routeConstructor,
		fx.As(new(R)),
		fx.ResultTags(string(acc)),
	)
}

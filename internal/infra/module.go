package infra

import "go.uber.org/fx"

var Module = fx.Module(
	"Infra",
	fx.Provide(NewDatabase),
	fx.Provide(NewAppLogger),
	fx.Provide(NewDBLogger),
)

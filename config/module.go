package config

import "go.uber.org/fx"

var Module = fx.Module(
	"Config",
	fx.Provide(NewEnv),
)

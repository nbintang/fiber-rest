package infra

import "go.uber.org/fx"

var Module = fx.Module(
	"infra",
	fx.Provide(
		NewDatabaseService,
		NewRedisService,
	),
	fx.Provide(
		NewValidatorService,
		NewTokenService,
		NewEmailService,
	),
	fx.Provide(
		NewLogger,
		NewDBLogger,
	),
	fx.Invoke(
		RegisterRedisLifecycle,
		RegisterDatabaseLifecycle,
	),
)

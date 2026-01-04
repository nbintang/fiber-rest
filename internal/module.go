package internal

import (
	"rest-fiber/internal/auth"
	"rest-fiber/internal/user"

	"go.uber.org/fx"
)

var RunApp = func(lc fx.Lifecycle, a *App) { a.Run(lc) }

var BusinessModules = fx.Module(
	"Modules",
	user.Module,
	auth.Module,
)

var Module = fx.Module(
	"App",
	BusinessModules,
	fx.Provide(NewApp),
	fx.Invoke(
		RegisterAllRoutes,
		RunApp,
	),
)

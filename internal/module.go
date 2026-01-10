package internal

import (
	"rest-fiber/internal/auth"
	"rest-fiber/internal/category"
	"rest-fiber/internal/post"
	"rest-fiber/internal/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"app",
	//Business Modules
	fx.Options(
		category.Module,
		post.Module,
		user.Module,
		auth.Module,
	),
	//Provide App
	fx.Provide(NewBootstrap),
	fx.Invoke(
		RegisterAllRoutes,
		RegisterHttpLifecycle,
	),
)

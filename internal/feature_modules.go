package internal

import (
	"rest-fiber/internal/auth"
	"rest-fiber/internal/category"
	"rest-fiber/internal/post"
	"rest-fiber/internal/user"

	"go.uber.org/fx"
)

var FeatureModules = fx.Options(
	category.Module,
	post.Module,
	user.Module,
	auth.Module,
)

package main

import (
	"rest-fiber/config"
	app "rest-fiber/internal"
	"rest-fiber/internal/infra"
	"rest-fiber/pkg"

	"go.uber.org/fx"
)


func main() {
	pkg.LoadEnv()
	fx.New(
		config.Module,
		infra.Module,
		app.Module,
	).Run()
}

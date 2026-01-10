package main

import (
	"rest-fiber/config"
	app "rest-fiber/internal"
	"rest-fiber/internal/infra"
	"rest-fiber/utils"

	"go.uber.org/fx"
)


func main() {
	utils.LoadEnv()
	fx.New(
		config.Module,
		infra.Module,
		app.Module,
	).Run()
}

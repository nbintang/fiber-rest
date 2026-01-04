package main

import (
	"rest-fiber/config"
	app "rest-fiber/internal" 
	"rest-fiber/internal/infra" 

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
	fx.New(
		config.Module,
		infra.Module,
		app.Module,
	).Run()
}

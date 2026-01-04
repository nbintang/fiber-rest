package internal

import (
	"context"
	"log"
	"rest-fiber/config"
	"rest-fiber/internal/setup"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type App struct {
	Fiber *fiber.App
	API   fiber.Router
	Env   config.Env
}

func NewApp(env config.Env) *App {
	f := fiber.New(fiber.Config{
		ErrorHandler: setup.DefaultErrorHandler,
	})

	api := f.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"Test": "tod",
		})
	})

	return &App{
		Fiber: f,
		API:   api,
		Env:   env,
	}
}

func (a *App) Run(lc fx.Lifecycle) {
	addr := a.Env.AppAddr
	if addr == "" {
		addr = ":8080"
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Fiber listening on %s", addr)
				if err := a.Fiber.Listen(addr); err != nil {
					log.Printf("Fiber stopped: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return a.Fiber.Shutdown()
		},
	})
}

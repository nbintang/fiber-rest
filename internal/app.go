package internal

import (
	"context"
	"log"
	"rest-fiber/config"
	"rest-fiber/internal/middleware"
	"rest-fiber/internal/setup"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type App struct {
	*fiber.App
	PublicRoute    fiber.Router
	ProtectedRoute fiber.Router
	Env            config.Env
}

func NewApp(env config.Env) *App {
	f := fiber.New(fiber.Config{
		ErrorHandler: setup.DefaultErrorHandler,
	})

	api := f.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Wellcome to API")
	})

	protected := api.Group("/protected");

	protected.Use(middleware.AuthAccess(env), middleware.CurrentAuthUser())

	return &App{
		App:            f,
		PublicRoute:    api,
		ProtectedRoute: protected,
		Env:            env,
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
				if err := a.Listen(addr); err != nil {
					log.Printf("Fiber stopped: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return a.Shutdown()
		},
	})
}

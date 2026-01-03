package internal

import (
	"context"
	"log"
	"rest-fiber/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func NewApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			statusCode := fiber.StatusInternalServerError
			msg := "Internal Server Error"
			if e, ok := err.(*fiber.Error); ok {
				statusCode = e.Code
				msg = e.Message
			}
			return c.Status(statusCode).JSON(fiber.Map{
				"error": msg,
				"status": statusCode,
				"timestamp": time.Now().Unix(),
			})
		},
	})
}

func RunApp(lc fx.Lifecycle, app *fiber.App, env config.Env) {
	addr := env.AppAddr
	if addr == "" {
		addr = ":3000"
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Fiber listening on %s", addr)
				if err := app.Listen(addr); err != nil {
					log.Printf("Fiber stopped: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}

var Module = fx.Module(
	"App",
	fx.Provide(NewApp),
	fx.Invoke(RunApp),
)

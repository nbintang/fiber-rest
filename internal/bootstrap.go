package internal

import (
	"rest-fiber/config"
	"rest-fiber/internal/infra"
	"rest-fiber/internal/middleware"
	"rest-fiber/internal/setup"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Bootstrap struct {
	*fiber.App
	PublicRoute    fiber.Router
	ProtectedRoute fiber.Router
	Env            config.Env
	Logger         *infra.AppLogger
}

func NewBootstrap(env config.Env, logger *infra.AppLogger, redisService infra.RedisService) *Bootstrap {
	app := fiber.New(fiber.Config{
		ErrorHandler: setup.DefaultErrorHandler,
		AppName:      "Fiber Rest API",
	})
	app.Use(setup.LoggerRequest(logger))
	app.Use(cors.New(cors.ConfigDefault))
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Wellcome to API")
	})

	protected := api.Group("/protected")
	// Protected Routes Provider
	protected.Use(
		middleware.AuthAccessToken(env),
		middleware.AccessCurrentUser(),
		middleware.AccessNotBlacklisted(redisService),
	)

	return &Bootstrap{
		App:            app,
		PublicRoute:    api,
		ProtectedRoute: protected,
		Env:            env,
		Logger:         logger,
	}
}

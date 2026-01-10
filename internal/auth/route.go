package auth

import (
	"rest-fiber/config"
	"rest-fiber/internal/contract"
	"rest-fiber/internal/infra"
	"rest-fiber/internal/middleware"
	"rest-fiber/pkg/httpx"

	"github.com/gofiber/fiber/v2" 
)

type AuthRouteParams struct {
	httpx.RouteParams
	AuthHandler AuthHandler
	Env         config.Env
}

type authRouteImpl struct {
	authHandler AuthHandler
	env         config.Env
}

func NewAuthRoute(params AuthRouteParams) contract.Route {
	return &authRouteImpl{authHandler: params.AuthHandler, env: params.Env}
}
func (r *authRouteImpl) RegisterRoute(api fiber.Router) {
	auth := api.Group("/auth")

	redisStarage := infra.GetRedisStorage(r.env)
	storageParams := middleware.RateLimiterParams{MaxLimit: 3, Storage: redisStarage}
	auth.Post("/register", middleware.AuthRateLimit(storageParams), r.authHandler.Register)
	auth.Post("/verify", middleware.AuthRateLimit(storageParams), r.authHandler.VerifyEmail)
	auth.Post("/login", middleware.AuthRateLimit(storageParams), r.authHandler.Login)
	auth.Delete("/logout", middleware.AuthRateLimit(storageParams), r.authHandler.Logout)
	auth.Post("/refresh-token", middleware.AuthRateLimit(storageParams), r.authHandler.RefreshToken)
}

package auth

import (
	"rest-fiber/config"
	"rest-fiber/internal/contract"
	"rest-fiber/internal/infra"
	"rest-fiber/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type AuthRouteParams struct {
	fx.In
	H   AuthHandler
	Env config.Env
}

type authRouteImpl struct {
	h   AuthHandler
	env config.Env
}

func NewAuthRoute(p AuthRouteParams) contract.Route {
	return &authRouteImpl{h: p.H, env: p.Env}
}
func (r *authRouteImpl) RegisterRoute(api fiber.Router) {
	auth := api.Group("/auth")

	redisStarage := infra.GetRedisStorage(r.env)
	storageParams := middleware.RateLimiterParams{MaxLimit: 3, Storage: redisStarage}
	auth.Post("/register", middleware.RateLimiter(storageParams), r.h.Register)
	auth.Post("/verify", middleware.RateLimiter(storageParams), r.h.VerifyEmail)
	auth.Post("/login", middleware.RateLimiter(storageParams), r.h.Login)
	auth.Delete("/logout", middleware.RateLimiter(storageParams), r.h.Logout)
	auth.Post("/refresh-token", middleware.RateLimiter(storageParams), r.h.RefreshToken)
}

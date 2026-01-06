package auth

import (
	"rest-fiber/config"
	"rest-fiber/internal/contract"
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

	auth.Post("/register", middleware.RateLimiter(middleware.RateLimiterParams{MaxLimit: 5}), r.h.Register)
	auth.Post("/verify", middleware.RateLimiter(middleware.RateLimiterParams{MaxLimit: 5}), r.h.VerifyEmail)
	auth.Post("/login", middleware.RateLimiter(middleware.RateLimiterParams{MaxLimit: 5}), r.h.Login)
	auth.Post("/refresh-token", r.h.RefreshToken)
}

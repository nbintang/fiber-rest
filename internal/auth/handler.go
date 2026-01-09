package auth

import (
	"rest-fiber/config"
	"rest-fiber/internal/infra" 
	"rest-fiber/pkg/httpx"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authHandlerImpl struct {
	authService AuthService
	validate    infra.Validator
	env         config.Env
	logger      *infra.AppLogger
}

func NewAuthHandler(authService AuthService, validate infra.Validator, env config.Env, logger *infra.AppLogger) AuthHandler {
	return &authHandlerImpl{authService, validate, env, logger}
}

func (h *authHandlerImpl) Register(c *fiber.Ctx) error {
	var dto RegisterRequestDTO
	ctx := c.UserContext()

	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(dto); err != nil {
		return err
	}

	if err := h.authService.Register(ctx, &dto); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(httpx.NewHttpResponse(fiber.StatusCreated, "Success! please check your email", nil))
}

func (h *authHandlerImpl) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	ctx := c.UserContext()
	tokens, err := h.authService.VerifyEmailToken(ctx, token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(
		fiber.StatusOK,
		"Email Verified Successfull!",
		fiber.Map{"access_token": token},
	))
}

func (h *authHandlerImpl) Login(c *fiber.Ctx) error {
	var dto LoginRequestDTO
	ctx := c.UserContext()
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	tokens, err := h.authService.Login(ctx, &dto)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	return c.Status(fiber.StatusOK).JSON(
		httpx.NewHttpResponse(
			fiber.StatusOK,
			"Login successful",
			fiber.Map{"access_token": tokens.AccessToken},
		))
}

func (h *authHandlerImpl) RefreshToken(c *fiber.Ctx) error {
	ctx := c.UserContext()
	oldRefreshToken := c.Cookies("refresh_token")
	if oldRefreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "no refresh token provided")
	}

	tokens, err := h.authService.RefreshToken(ctx, oldRefreshToken)
	if err != nil {
		h.clearRefreshTokenCookie(c)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(
		fiber.StatusOK,
		"Refresh Token successful",
		fiber.Map{"access_token": tokens.AccessToken},
	))
}

func (h *authHandlerImpl) Logout(c *fiber.Ctx) error {
	ctx := c.UserContext()
	refreshToken := c.Cookies("refresh_token")
	err := h.authService.Logout(ctx, refreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no token provide, please login!")
	}
	h.clearRefreshTokenCookie(c)
	return c.Status(fiber.StatusBadRequest).JSON(httpx.NewHttpResponse(fiber.StatusOK, "logout Success", nil))
}

func (h *authHandlerImpl) setRefreshTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
	})
}
func (h *authHandlerImpl) clearRefreshTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
	})
}

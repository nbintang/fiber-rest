package auth

import (
	"rest-fiber/config"
	"rest-fiber/internal/infra"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	authService AuthService
	validate    infra.Validator
	env         config.Env
}

func NewAuthHandler(authService AuthService, validate infra.Validator, env config.Env) AuthHandler {
	return &authHandler{authService, validate, env}
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	var dto RegisterRequestDTO
	ctx := c.UserContext()

	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(dto); err != nil {
		return err
	}

	if err := h.authService.Register(ctx, &dto); err != nil {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "token sent successfully, please check your email",
	})
}

func (h *authHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	ctx := c.UserContext()
	token, err := h.authService.VerifyEmailToken(ctx, token)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(200).JSON(fiber.Map{
		"message":      "email verified successfully",
		"access_token": token,
	})
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var dto LoginRequestDTO
	ctx := c.UserContext()

	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	token, err := h.authService.Login(ctx, &dto)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(200).JSON(fiber.Map{
		"message":      "login successful",
		"access_token": token,
	})
}

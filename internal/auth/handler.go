package auth

import (
	"errors"
	"rest-fiber/pkg"

	"github.com/gofiber/fiber/v2"
)


type authHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) AuthHandler {
	return &authHandler{authService}
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	var dto RegisterRequestDTO
	ctx := c.UserContext()
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := pkg.Validate.Struct(dto); err != nil {
		return err
	}
	if err := h.authService.Register(ctx, &dto); err != nil {
		if errors.Is(err, pkg.ErrAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, "User already exists")
		}
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfuly",
	})
}

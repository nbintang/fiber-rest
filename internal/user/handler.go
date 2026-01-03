package user

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userResponses, err := h.userService.FindAllUsers(ctx)
	if err != nil {
		return err
	}
	return c.JSON(userResponses)
}

func (h *userHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()

	userResponse, err := h.userService.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("Not Found")) {
			return fiber.NewError(fiber.StatusNotFound, "User Not Found")
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

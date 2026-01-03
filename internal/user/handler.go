package user

import (
	"errors"
	"rest-fiber/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}
	ctx := c.UserContext()

	userResponse, err := h.userService.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "User Not Found")
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

package user

import (
	"rest-fiber/internal/infra"
	"rest-fiber/internal/setup"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userHandlerImpl struct {
	userService UserService
	logger      *infra.AppLogger
}

func NewUserHandler(userService UserService, logger *infra.AppLogger) UserHandler {
	return &userHandlerImpl{userService, logger}
}

func (h *userHandlerImpl) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userResponses, err := h.userService.FindAllUsers(ctx)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(setup.NewHttpResponse(fiber.StatusOK, "Success", userResponses))
}

func (h *userHandlerImpl) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}
	ctx := c.UserContext()

	userResponse, err := h.userService.FindUserByID(ctx, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(setup.NewHttpResponse(fiber.StatusOK, "Success", userResponse))
}

func (h *userHandlerImpl) GetCurrentUser(c *fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	if userId == "" {
		c.Status(401).JSON(setup.NewHttpResponse(fiber.StatusUnauthorized, "Unauthorized", nil))
	}
	ctx := c.UserContext()
	h.logger.Infof("user Id :%s", userId)
	userResponse, err := h.userService.FindUserByID(ctx, userId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(setup.NewHttpResponse(fiber.StatusOK, "Success", userResponse))
}

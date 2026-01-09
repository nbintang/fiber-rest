package user

import (
	"rest-fiber/internal/infra"
	"rest-fiber/pkg/httpx"

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
	var query httpx.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := query.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	query = query.Normalize(10, 100)

	data, total, err := h.userService.FindAllUsers(ctx, query.Page, query.Limit, query.Offset())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	meta := httpx.NewPaginationResponse(query.Page, query.Limit, total)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    data,
		"meta":    meta,
	})
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

	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(fiber.StatusOK, "Success", userResponse))
}

func (h *userHandlerImpl) GetCurrentUser(c *fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	if userId == "" {
		c.Status(401).JSON(httpx.NewHttpResponse(fiber.StatusUnauthorized, "Unauthorized", nil))
	}
	ctx := c.UserContext()
	h.logger.Infof("user Id :%s", userId)
	userResponse, err := h.userService.FindUserByID(ctx, userId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(fiber.StatusOK, "Success", userResponse))
}

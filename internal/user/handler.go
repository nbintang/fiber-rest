package user

import (
	"rest-fiber/internal/infra"
	"rest-fiber/internal/middleware"
	"rest-fiber/pkg/httpx"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userHandlerImpl struct {
	userService UserService
	logger      *infra.AppLogger
	validator   infra.Validator
}

func NewUserHandler(userService UserService, logger *infra.AppLogger, validator infra.Validator) UserHandler {
	return &userHandlerImpl{userService, logger, validator}
}

func (h *userHandlerImpl) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var query httpx.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	query = query.Normalize(10, 100)

	data, total, err := h.userService.FindAllUsers(ctx, query.Page, query.Limit, query.Offset())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	meta := httpx.NewPaginationMeta(query.Page, query.Limit, total)
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpPaginationResponse[[]UserResponseDTO](
		fiber.StatusOK,
		"Success",
		data,
		meta,
	))
}

func (h *userHandlerImpl) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}
	ctx := c.UserContext()
	data, err := h.userService.FindUserByID(ctx, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(fiber.StatusOK, "Success", data))
}

func (h *userHandlerImpl) GetCurrentUserProfile(c *fiber.Ctx) error {
	currentUser, err := middleware.GetCurrentUser(c)
	if err != nil {
		return err
	}
	ctx := c.UserContext()
	h.logger.Infof("user Id :%s", currentUser.ID)
	data, err := h.userService.FindUserByID(ctx, currentUser.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(fiber.StatusOK, "Success", data))
}

func (h *userHandlerImpl) UpdateCurrentUser(c *fiber.Ctx) error {
	currentUser, err := middleware.GetCurrentUser(c)
	if err != nil {
		return err
	}
	var body UserUpdateDTO
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.validator.Struct(body); err != nil {
		return err
	}

	if err := h.userService.UpdateProfile(c.UserContext(), currentUser.ID, body); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(
		httpx.NewHttpResponse[any](fiber.StatusOK, "User Updated Successfully", nil),
	)
}

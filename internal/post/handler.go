package post

import (
	"rest-fiber/internal/identity"
	"rest-fiber/internal/infra/validator"
	"rest-fiber/pkg/httpx"
	"rest-fiber/pkg/pagination"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type postHandlerImpl struct {
	postService PostService
	validate    validator.Service
}

func NewPostHandler(postService PostService, validate validator.Service) PostHandler {
	return &postHandlerImpl{postService, validate}
}

func (h *postHandlerImpl) GetAllPosts(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var query pagination.Query
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	query = query.Normalize(10, 100)

	data, total, err := h.postService.FindAllPosts(ctx, query.Page, query.Limit, query.Offset())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	meta := pagination.NewPaginationMeta(query.Page, query.Limit, total)
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpPaginationResponse[[]PaginatedPostResponseDTO](
		fiber.StatusOK,
		"Success",
		data,
		meta,
	))
}
func (h *postHandlerImpl) GetPostByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}
	ctx := c.UserContext()
	data, err := h.postService.FindPostByID(ctx, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(fiber.StatusOK, "Success", data))
}
func (h *postHandlerImpl) CreatePost(c *fiber.Ctx) error {
	currentUser, err := identity.CurrentUser(c)
	if err != nil {
		return err
	}

	var body PostRequestDTO
	ctx := c.UserContext()
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(body); err != nil {
		return err
	}

	data, err := h.postService.CreatePost(ctx, body, currentUser.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(httpx.NewHttpResponse[*PostResponseDTO](
		fiber.StatusCreated,
		"Success",
		data,
	))
}
func (h *postHandlerImpl) UpdatePostByID(c *fiber.Ctx) error {
	currentUser, err := identity.CurrentUser(c)
	if err != nil {
		return err
	}

	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	var body PostRequestDTO
	ctx := c.UserContext()
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(body); err != nil {
		return err
	}

	data, err := h.postService.UpdatePostByID(ctx, id, body, currentUser.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse[*PostResponseDTO](
		fiber.StatusOK,
		"Success",
		data,
	))
}

func (h *postHandlerImpl) DeletePostByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}
	ctx := c.UserContext()
	if err := h.postService.DeletePostByID(ctx, id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse[any](
		fiber.StatusOK,
		"Deleted Successfully	",
		nil,
	))
}

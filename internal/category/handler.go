package category

import ( 
	"rest-fiber/internal/infra/validator"
	"rest-fiber/pkg/httpx"
	"rest-fiber/pkg/pagination"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type categoryHandlerImpl struct {
	categoryService CategoryService
	validate        validator.Service
}

func NewCategoryHandler(categoryService CategoryService, validate validator.Service) CategoryHandler {
	return &categoryHandlerImpl{categoryService, validate}
}

func (h *categoryHandlerImpl) GetAllCategories(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var query pagination.Query
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	query = query.Normalize(10, 100)

	data, total, err := h.categoryService.FindAllCategories(ctx, query.Page, query.Limit, query.Offset())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	meta := pagination.NewPaginationMeta(query.Page, query.Limit, total)
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpPaginationResponse[[]CategoryResponseDTO](
		fiber.StatusOK,
		"Success",
		data,
		meta,
	))
}

func (h *categoryHandlerImpl) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}
	ctx := c.UserContext()
	data, err := h.categoryService.FindCategoryByID(ctx, id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(
		fiber.StatusOK,
		"Success",
		data,
	))
}
func (h *categoryHandlerImpl) CreateCategory(c *fiber.Ctx) error {
	var body CategoryRequestDTO
	ctx := c.UserContext()
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	data, err := h.categoryService.CreateCategory(ctx, &body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(
		fiber.StatusOK,
		"Success",
		data,
	))
}
func (h *categoryHandlerImpl) UpdateCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	var body CategoryRequestDTO
	ctx := c.UserContext()
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	data, err := h.categoryService.UpdateCategoryByID(ctx, id, &body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(
		fiber.StatusOK,
		"Success",
		data,
	))
}
func (h *categoryHandlerImpl) DeleteCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}
	ctx := c.Context()
	if err := h.categoryService.DeleteCategoryByID(ctx, id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse[any](
		fiber.StatusOK,
		"Deleted Successfully",
		nil,
	))
}

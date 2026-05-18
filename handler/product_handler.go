package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/service"
)

type ProductHandler struct {
	Service service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{Service: svc}
}

// Search godoc
// @Summary Search products
// @Router /products/search [post]
func (h *ProductHandler) Search(c *fiber.Ctx) error {
	var req model.ProductSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: "invalid request body"})
	}

	result, err := h.Service.Search(c.UserContext(), req)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Message: "success", Data: result})
}

// GetByID godoc
// @Summary Get product by ID
// @Router /products/:id [get]
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: "invalid id"})
	}

	product, err := h.Service.GetByID(c.UserContext(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Message: "success", Data: product})
}

// Create godoc
// @Summary Create product
// @Router /products [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req model.ProductCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: "invalid request body"})
	}

	// ในของจริงดึงจาก middleware auth
	userID := c.Get("X-User-ID", "anonymous")

	product, err := h.Service.Create(c.UserContext(), req, userID)
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(model.Response{Message: "created", Data: product})
}

// Update godoc
// @Summary Update product
// @Router /products/:id [put]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: "invalid id"})
	}

	var req model.ProductUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: "invalid request body"})
	}

	userID := c.Get("X-User-ID", "anonymous")

	product, err := h.Service.Update(c.UserContext(), id, req, userID)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Message: "updated", Data: product})
}

// Delete godoc
// @Summary Delete product
// @Router /products/:id [delete]
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: "invalid id"})
	}

	if err := h.Service.Delete(c.UserContext(), id); err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Message: "deleted"})
}

// ─── Error Handler ─────────────────────────────────────────

func handleError(c *fiber.Ctx, err error) error {
	var notFound *model.NotFoundError
	if errors.As(err, &notFound) {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{Message: "not found"})
	}

	var valErr *model.ValidationError
	if errors.As(err, &valErr) {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: valErr.Message})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(model.Response{Message: "internal server error"})
}

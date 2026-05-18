package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/service"
	"github.com/iiimomoniii/inventory_backend/utils"
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
		return utils.InvalidBody(c)
	}
	result, err := h.Service.Search(c.UserContext(), req)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Status: fiber.StatusOK, Message: "success", Data: result})
}

// GetByID godoc
// @Summary Get product by ID
// @Router /products/:id [get]
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "PRD005")
	}
	product, err := h.Service.GetByID(c.UserContext(), id)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Status: fiber.StatusOK, Message: "success", Data: product})
}

// Create godoc — Single
// @Summary Create product (single)
// @Router /products [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req model.ProductCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.InvalidBody(c)
	}

	userID := c.Get("X-User-ID", "anonymous")

	product, err := h.Service.Create(c.UserContext(), req, userID)
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Status:  fiber.StatusCreated,
		Message: "created",
		Data:    product,
	})
}

// ItemsCreate godoc — Items
// @Summary Create products (Items) — error ของใครของมัน
// @Router /products/Items [post]
func (h *ProductHandler) CreateItems(c *fiber.Ctx) error {
	var items []model.ProductCreateRequest
	if err := c.BodyParser(&items); err != nil {
		return utils.InvalidBody(c)
	}
	if len(items) == 0 {
		return utils.InvalidBody(c)
	}

	userID := c.Get("X-User-ID", "anonymous")
	results := make([]model.ItemResult, 0, len(items))

	for i, req := range items {
		// ลอง create แต่ละรายการ ไม่หยุดเมื่อ error
		product, err := h.Service.Create(c.UserContext(), req, userID)
		if err != nil {
			var valErr *model.ValidationError
			if errors.As(err, &valErr) {
				// error ของ item นี้ → บันทึก แล้วทำต่อ
				results = append(results, utils.BuildErrorItem(i, valErr.Code))
				continue
			}
			results = append(results, utils.BuildErrorItem(i, "GLB002"))
			continue
		}
		// สำเร็จ
		results = append(results, utils.BuildSuccessItem(i, product))
	}

	// ส่ง Items response — แสดง error ของใครของมัน
	return utils.ItemsResponse(c, results)
}

// Update godoc
// @Summary Update product
// @Router /products/:id [put]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "PRD005")
	}
	var req model.ProductUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.InvalidBody(c)
	}
	userID := c.Get("X-User-ID", "anonymous")
	product, err := h.Service.Update(c.UserContext(), id, req, userID)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Status: fiber.StatusOK, Message: "updated", Data: product})
}

// Delete godoc
// @Summary Delete product
// @Router /products/:id [delete]
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "PRD005")
	}
	if err := h.Service.Delete(c.UserContext(), id); err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Status: fiber.StatusOK, Message: "deleted"})
}

// ─── Error Handler ─────────────────────────────────────────

func handleError(c *fiber.Ctx, err error) error {
	var notFound *model.NotFoundError
	if errors.As(err, &notFound) {
		return utils.NotFound(c)
	}
	var valErr *model.ValidationError
	if errors.As(err, &valErr) {
		return utils.BadRequest(c, valErr.Code)
	}
	return utils.InternalError(c)
}

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

// getUserID — ดึง username จาก JWT locals ใช้เป็น created_by / updated_by
func getUserID(c *fiber.Ctx) string {
	if userID, ok := c.Locals("userID").(string); ok && userID != "" {
		return userID
	}
	return "anonymous"
}

// Search godoc
// @Summary Search products
// @Router /v1/products/search [post]
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
// @Router /v1/products/:id [get]
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

// Create godoc
// @Summary Create product (single)
// @Router /v1/products/create [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req model.ProductCreateRequest
	if err := utils.StrictUnmarshal(c.Body(), &req); err != nil {
		return utils.InvalidBody(c)
	}

	product, err := h.Service.Create(c.UserContext(), req, getUserID(c))
	if err != nil {
		var serviceErr *service.ServiceError
		if errors.As(err, &serviceErr) && serviceErr.Code == "PRD010" {
			return utils.Conflict(c, serviceErr.Code)
		}

		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Status:  fiber.StatusCreated,
		Message: "created",
		Data:    product,
	})
}

// CreateItems godoc
// @Summary Create products (bulk)
// @Router /v1/products/create/items [post]
func (h *ProductHandler) CreateItems(c *fiber.Ctx) error {
	var items []model.ProductCreateRequest
	if err := c.BodyParser(&items); err != nil {
		return utils.InvalidBody(c)
	}
	if len(items) == 0 {
		return utils.InvalidBody(c)
	}

	// validate แต่ละ item ว่าเป็น camelCase
	for i := range items {
		if err := utils.StrictJSONKeys(c.Body(), model.ProductCreateRequest{}); err != nil {
			return utils.InvalidBody(c)
		}
		_ = i
	}

	userID := getUserID(c)
	results := make([]model.ItemResult, 0, len(items))

	for i, req := range items {
		product, err := h.Service.Create(c.UserContext(), req, userID)
		if err != nil {
			var valErr *model.ValidationError
			if errors.As(err, &valErr) {
				results = append(results, utils.BuildErrorItem(i, valErr.Code))
				continue
			}
			results = append(results, utils.BuildErrorItem(i, "GLB002"))
			continue
		}
		results = append(results, utils.BuildSuccessItem(i, product))
	}
	return utils.ItemsResponse(c, results)
}

// Update godoc
// @Summary Update product
// @Router /v1/products/update/:id [put]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "PRD005")
	}
	var req model.ProductUpdateRequest
	if err := utils.StrictUnmarshal(c.Body(), &req); err != nil {
		return utils.InvalidBody(c)
	}
	product, err := h.Service.Update(c.UserContext(), id, req, getUserID(c))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(model.Response{Status: fiber.StatusOK, Message: "updated", Data: product})
}

// Delete godoc
// @Summary Delete product
// @Router /v1/products/delete/:id [delete]
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

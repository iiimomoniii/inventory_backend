package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/repository"
	"github.com/iiimomoniii/inventory_backend/utils"
)

type CategoryHandler struct {
	Repo repository.CategoryRepository
}

func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{Repo: repo}
}

// GetAll godoc
// @Summary Get all categories
// @Router /v1/categories [get]
func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	categories, err := h.Repo.FindAll(c.UserContext())
	if err != nil {
		return utils.InternalError(c)
	}
	return c.JSON(model.Response{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    categories,
	})
}

// GetByID godoc
// @Summary Get category by ID
// @Router /v1/categories/:id [get]
func (h *CategoryHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.BadRequest(c, "PRD005")
	}

	category, err := h.Repo.FindByID(c.UserContext(), id)
	if err != nil {
		var notFound *model.NotFoundError
		if errors.As(err, &notFound) {
			return utils.NotFound(c)
		}
		return utils.InternalError(c)
	}
	return c.JSON(model.Response{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    category,
	})
}

package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/service"
	"github.com/iiimomoniii/inventory_backend/utils"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{Service: svc}
}

// Create godoc
// @Summary Create user
// @Router /v1/users/create [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req model.UserCreateRequest
	if err := utils.StrictUnmarshal(c.Body(), &req); err != nil {
		return utils.CustomErrorResp(fiber.StatusBadRequest, "GLB003", c)
	}

	createdBy := strings.TrimSpace(req.Username)

	user, err := h.Service.Create(c.UserContext(), req, createdBy)
	if err != nil {
		var svcErr *service.ServiceError
		if errors.As(err, &svcErr) {
			return utils.CustomErrorResp(fiber.StatusBadRequest, svcErr.Code, c)
		}
		return utils.CustomErrorResp(fiber.StatusInternalServerError, "GLB002", c)
	}

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Status:  fiber.StatusCreated,
		Message: "created",
		Data:    user,
	})
}

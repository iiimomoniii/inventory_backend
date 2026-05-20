package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/service"
	"github.com/iiimomoniii/inventory_backend/utils"
)

type AuthHandler struct {
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// GenerateToken godoc
// @Summary Generate JWT token
// @Router /auth/token [post]
func (h *AuthHandler) GenerateToken(c *fiber.Ctx) error {
	var req model.TokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.InvalidBody(c)
	}

	if req.Username == "" || req.Password == "" {
		return utils.InvalidBody(c)
	}

	tokenResp, err := h.AuthService.Login(c.UserContext(), req.Username, req.Password)
	if err != nil {
		var valErr *model.ValidationError
		if errors.As(err, &valErr) {
			return utils.Unauthorized(c)
		}
		return utils.InternalError(c)
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    tokenResp,
	})
}

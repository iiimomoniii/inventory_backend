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
		return utils.CustomErrorResp(fiber.StatusBadRequest, "GLB003", c)
	}
	if req.Username == "" || req.Password == "" {
		return utils.CustomErrorResp(fiber.StatusBadRequest, "GLB003", c)
	}

	tokenResp, err := h.AuthService.Login(c.UserContext(), req.Username, req.Password)
	if err != nil {
		var svcErr *service.ServiceError
		if errors.As(err, &svcErr) {
			return utils.CustomErrorResp(fiber.StatusUnauthorized, svcErr.Code, c)
		}
		return utils.CustomErrorResp(fiber.StatusInternalServerError, "GLB002", c)
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    tokenResp,
	})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req model.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.CustomErrorResp(fiber.StatusBadRequest, "GLB003", c)
	}
	if req.RefreshToken == "" {
		return utils.CustomErrorResp(fiber.StatusBadRequest, "GLB003", c)
	}

	tokenResp, err := h.AuthService.Refresh(c.UserContext(), req.RefreshToken)
	if err != nil {
		var svcErr *service.ServiceError
		if errors.As(err, &svcErr) {
			return utils.CustomErrorResp(fiber.StatusUnauthorized, svcErr.Code, c)
		}
		return utils.CustomErrorResp(fiber.StatusInternalServerError, "GLB002", c)
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    tokenResp,
	})
}

// Logout godoc
// @Summary Logout
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req model.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.CustomErrorResp(fiber.StatusBadRequest, "GLB003", c)
	}
	if req.RefreshToken == "" {
		return utils.CustomErrorResp(fiber.StatusBadRequest, "GLB003", c)
	}

	if err := h.AuthService.Logout(c.UserContext(), req.RefreshToken); err != nil {
		var svcErr *service.ServiceError
		if errors.As(err, &svcErr) {
			return utils.CustomErrorResp(fiber.StatusUnauthorized, svcErr.Code, c)
		}
		return utils.CustomErrorResp(fiber.StatusInternalServerError, "GLB002", c)
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Status:  fiber.StatusOK,
		Message: "logout success",
	})
}

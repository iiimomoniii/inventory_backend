package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/iiimomoniii/inventory_backend/model"
	"github.com/iiimomoniii/inventory_backend/service"
	"github.com/iiimomoniii/inventory_backend/utils"
)

func handleError(c *fiber.Ctx, err error) error {

	var notFound *model.NotFoundError
	if errors.As(err, &notFound) {
		return utils.NotFound(c)
	}

	var valErr *model.ValidationError
	if errors.As(err, &valErr) {
		return utils.BadRequest(c, valErr.Code)
	}

	var serviceErr *service.ServiceError
	if errors.As(err, &serviceErr) {

		return utils.BadRequest(c, serviceErr.Code)
	}

	return utils.InternalError(c)
}

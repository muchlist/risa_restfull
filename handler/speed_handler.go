package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/risa_restfull/service"
)

func NewSpeedHandler(speedService service.SpeedTestServiceAssumer) *speedHandler {
	return &speedHandler{
		service: speedService,
	}
}

type speedHandler struct {
	service service.SpeedTestServiceAssumer
}

func (sh *speedHandler) Retrieve(c *fiber.Ctx) error {

	speedList, apiErr := sh.service.RetrieveSpeed(c.Context())
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": speedList})
}

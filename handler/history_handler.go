package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

func NewHistoryHandler(histService service.HistoryServiceAssumer) *historyHandler {
	return &historyHandler{
		service: histService,
	}
}

type historyHandler struct {
	service service.HistoryServiceAssumer
}

func (h *historyHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.HistoryRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	insertID, apiErr := h.service.InsertHistory(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan history berhasi, ID: %s", *insertID)}
	return c.JSON(res)
}

package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

func NewStockHandler(stockService service.StockServiceAssumer) *stockHandler {
	return &stockHandler{
		service: stockService,
	}
}

type stockHandler struct {
	service service.StockServiceAssumer
}

func (s *stockHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.StockRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	insertID, apiErr := s.service.InsertStock(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan stock berhasil, ID: %s", *insertID)}
	return c.JSON(res)
}

func (s *stockHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	stockID := c.Params("id")

	var req dto.StockEditRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	stockEdited, apiErr := s.service.EditStock(*claims, stockID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(stockEdited)
}

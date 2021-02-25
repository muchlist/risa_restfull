package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

func NewCctvHandler(cctvService service.CctvServiceAssumer) *cctvHandler {
	return &cctvHandler{
		service: cctvService,
	}
}

type cctvHandler struct {
	service service.CctvServiceAssumer
}

func (x *cctvHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.CctvRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	insertID, apiErr := x.service.InsertCctv(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan cctv berhasi, ID: %s", *insertID)}
	return c.JSON(res)
}

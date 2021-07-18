package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

func NewVendorCheckHandler(vendorCheckService service.VendorCheckServiceAssumer) *vendorCheckHandler {
	return &vendorCheckHandler{
		service: vendorCheckService,
	}
}

type vendorCheckHandler struct {
	service service.VendorCheckServiceAssumer
}

func (vc *vendorCheckHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	isVirtualCheck := len(c.Query("virtual")) != 0

	insertID, apiErr := vc.service.InsertVendorCheck(*claims, isVirtualCheck)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan cctv check berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

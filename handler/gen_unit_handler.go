package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

func NewGenUnitHandler(genUnitService service.GenUnitServiceAssumer) *genUnitHandler {
	return &genUnitHandler{
		service: genUnitService,
	}
}

type genUnitHandler struct {
	service service.GenUnitServiceAssumer
}

//Find menampilkan list unit
func (u *genUnitHandler) Find(c *fiber.Ctx) error {

	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	nameSearch := c.Query("name")
	categorySearch := c.Query("category")
	ipSearch := c.Query("ip")

	payload := dto.GenUnitFilter{
		Branch:   claims.Branch,
		Name:     nameSearch,
		Category: categorySearch,
		IP:       ipSearch,
	}

	userList, apiErr := u.service.FindUnit(payload)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"units": userList})
}

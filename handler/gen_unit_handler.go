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

func NewGenUnitHandler(genUnitService service.GenUnitServiceAssumer) *genUnitHandler {
	return &genUnitHandler{
		service: genUnitService,
	}
}

type genUnitHandler struct {
	service service.GenUnitServiceAssumer
}

//Find menampilkan list unit. Query name, category, ip
func (u *genUnitHandler) Find(c *fiber.Ctx) error {

	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	nameSearch := c.Query("name")
	categorySearch := c.Query("category")
	ipSearch := c.Query("ip")
	var pingsRetrieve bool
	if c.Query("pings") != "" {
		pingsRetrieve = true
	}

	payload := dto.GenUnitFilter{
		Branch:   claims.Branch,
		Name:     nameSearch,
		Category: categorySearch,
		IP:       ipSearch,
		Pings:    pingsRetrieve,
	}

	userList, apiErr := u.service.FindUnit(payload)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"units": userList})
}

//Find menampilkan list ip address. Query branch, category
func (u *genUnitHandler) GetIPList(c *fiber.Ctx) error {

	branch := c.Query("branch")
	category := c.Query("category")

	ipList, apiErr := u.service.GetIPList(branch, category)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"ip_list": ipList})
}

func (u genUnitHandler) UpdatePingState(c *fiber.Ctx) error {
	var req dto.GenUnitPingStateRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", c.IP(), err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", c.IP(), err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	count, apiErr := u.service.AppendPingState(req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("%d ip diupdate", count)}
	return c.JSON(res)
}

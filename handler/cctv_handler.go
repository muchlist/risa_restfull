package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils"
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
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	insertID, apiErr := x.service.InsertCctv(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan cctv berhasi, ID: %s", *insertID)}
	return c.JSON(res)
}

// GetCctv menampilkan cctvDetail
func (x *cctvHandler) GetCctv(c *fiber.Ctx) error {

	cctvID := c.Params("id")

	cctv, apiErr := x.service.GetCctvByID(cctvID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(cctv)
}

// Find menampilkan list cctv
// Query [branch, name, ip, location, disable]
func (x *cctvHandler) Find(c *fiber.Ctx) error {

	branch := c.Query("branch")
	name := c.Query("name")
	ip := c.Query("ip")
	location := c.Query("location")
	var disable bool
	if c.Query("disable") != "" {
		disable = true
	}

	filterA := dto.FilterBranchLocIPNameDisable{
		Branch:   branch,
		Location: location,
		IP:       ip,
		Name:     name,
		Disable:  disable,
	}

	cctvList, apiErr := x.service.FindCctv(filterA)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"cctv_list": cctvList})
}

//DisableCctv menghilangkan cctv dari list
// Param status [enable, disable]
func (x *cctvHandler) DisableCctv(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	userID := c.Params("id")
	status := c.Params("status")

	// validation
	statusAvailable := []string{"disable", "enable"}
	if !utils.InSlice(status, statusAvailable) {
		apiErr := rest_err.NewBadRequestError(fmt.Sprintf("Status yang dimasukkan tidak tersedia. gunakan %s", statusAvailable))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	var statusBool bool
	if status == "disable" {
		statusBool = true
	}

	cctvList, apiErr := x.service.DisableCctv(userID, *claims, statusBool)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"cctv_list": cctvList})
}

func (x *cctvHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := x.service.DeleteCctv(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": fmt.Sprintf("cctv %s berhasil dihapus", id)})
}

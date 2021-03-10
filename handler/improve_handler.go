package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
)

func NewImproveHandler(improveService service.ImproveServiceAssumer) *improveHandler {
	return &improveHandler{
		service: improveService,
	}
}

type improveHandler struct {
	service service.ImproveServiceAssumer
}

func (s *improveHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.ImproveRequest
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

	insertID, apiErr := s.service.InsertImprove(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan improvement berhasil, ID: %s", *insertID)}
	return c.JSON(res)
}

func (s *improveHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	improveID := c.Params("id")

	var req dto.ImproveEditRequest
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

	improveEdited, apiErr := s.service.EditImprove(*claims, improveID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(improveEdited)
}

func (s *improveHandler) ChangeImprove(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	improveID := c.Params("id")

	var req dto.ImproveChangeRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	improveEdited, apiErr := s.service.ChangeImprove(*claims, improveID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(improveEdited)
}

// GetImprove menampilkan improve Detail
func (s *improveHandler) GetImprove(c *fiber.Ctx) error {
	improveID := c.Params("id")

	improve, apiErr := s.service.GetImproveByID(improveID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(improve)
}

// Find menampilkan list improve
// Query [branch, c_status, start, end, limit]
func (s *improveHandler) Find(c *fiber.Ctx) error {
	branch := c.Query("branch")
	cStatus := stringToInt(c.Query("c_status"))
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	// limit default 300,
	filter := dto.FilterBranchCompleteTimeRangeLimit{
		FilterBranch:         branch,
		FilterCompleteStatus: cStatus,
		FilterStart:          int64(start),
		FilterEnd:            int64(end),
		Limit:                int64(limit),
	}

	improveList, apiErr := s.service.FindImprove(filter)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"improve_list": improveList})
}

// ActivateImprove mengaktifkan improve yang dibuat user selain yang memiliki hak
// Param status [enable, disable]
func (s *improveHandler) ActivateImprove(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	userID := c.Params("id")

	status := c.Params("status")

	// validation
	statusAvailable := []string{"disable", "enable"}
	if !sfunc.InSlice(status, statusAvailable) {
		apiErr := rest_err.NewBadRequestError(fmt.Sprintf("Status yang dimasukkan tidak tersedia. gunakan %s", statusAvailable))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	var isEnable bool
	if status == "enable" {
		isEnable = true
	}

	improveList, apiErr := s.service.ActivateImprove(userID, *claims, isEnable)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"improve_list": improveList})
}

func (s *improveHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := s.service.DeleteImprove(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": fmt.Sprintf("improve %s berhasil dihapus", id)})
}

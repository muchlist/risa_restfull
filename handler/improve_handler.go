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

func (iv *improveHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.ImproveRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := iv.service.InsertImprove(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan improvement berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (iv *improveHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	improveID := c.Params("id")

	var req dto.ImproveEditRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	improveEdited, apiErr := iv.service.EditImprove(c.Context(), *claims, improveID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": improveEdited})
}

func (iv *improveHandler) ChangeImprove(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	improveID := c.Params("id")

	var req dto.ImproveChangeRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	improveEdited, apiErr := iv.service.ChangeImprove(c.Context(), *claims, improveID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": improveEdited})
}

// GetImprove menampilkan improve Detail
func (iv *improveHandler) GetImprove(c *fiber.Ctx) error {
	improveID := c.Params("id")

	improve, apiErr := iv.service.GetImproveByID(c.Context(), improveID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": improve})
}

// Find menampilkan list improve
// Query [branch, c_status, start, end, limit]
func (iv *improveHandler) Find(c *fiber.Ctx) error {
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

	improveList, apiErr := iv.service.FindImprove(c.Context(), filter)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": improveList})
}

// ActivateImprove mengaktifkan improve yang dibuat user selain yang memiliki hak
// Param status [enable, disable]
func (iv *improveHandler) ActivateImprove(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	userID := c.Params("id")

	status := c.Params("status")

	// validation
	statusAvailable := []string{"disable", "enable"}
	if !sfunc.InSlice(status, statusAvailable) {
		apiErr := rest_err.NewBadRequestError(fmt.Sprintf("Status yang dimasukkan tidak tersedia. gunakan %s", statusAvailable))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	var isEnable bool
	if status == "enable" {
		isEnable = true
	}

	improveList, apiErr := iv.service.ActivateImprove(c.Context(), userID, *claims, isEnable)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": improveList})
}

func (iv *improveHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := iv.service.DeleteImprove(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("improve %s berhasil dihapus", id)})
}

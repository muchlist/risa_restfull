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

func NewAltaiPhyCheckHandler(altaiPhyCheckService service.AltaiPhyCheckServiceAssumer) *altaiPhyCheckHandler {
	return &altaiPhyCheckHandler{
		service: altaiPhyCheckService,
	}
}

type altaiPhyCheckHandler struct {
	service service.AltaiPhyCheckServiceAssumer
}

func (vc *altaiPhyCheckHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := vc.service.InsertAltaiPhyCheck(c.Context(), *claims, req.Name, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": insertID})
}

func (vc *altaiPhyCheckHandler) InsertQuarter(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	insertID, apiErr := vc.service.InsertAltaiPhyCheck(c.Context(), *claims, req.Name, true)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": insertID})
}

func (vc *altaiPhyCheckHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := vc.service.DeleteAltaiPhyCheck(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("altai check %s berhasil dihapus", id)})
}

func (vc *altaiPhyCheckHandler) Get(c *fiber.Ctx) error {
	checkID := c.Params("id")

	check, apiErr := vc.service.GetAltaiPhyCheckByID(c.Context(), checkID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": check})
}

// Find menampilkan list check
// Query [branch, start, end, limit]
func (vc *altaiPhyCheckHandler) Find(c *fiber.Ctx) error {
	branch := c.Query("branch")
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	checkList, apiErr := vc.service.FindAltaiPhyCheck(c.Context(), branch, dto.FilterTimeRangeLimit{
		FilterStart: int64(start),
		FilterEnd:   int64(end),
		Limit:       int64(limit),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": checkList})
}

func (vc *altaiPhyCheckHandler) UpdateCheckItem(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.AltaiPhyCheckItemUpdateRequest
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

	checkUpdated, apiErr := vc.service.UpdateAltaiPhyCheckItem(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": checkUpdated})
}

func (vc *altaiPhyCheckHandler) BulkUpdateCheckItem(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var reqs dto.BulkAltaiPhyCheckUpdateRequest
	if err := c.BodyParser(&reqs); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	for _, req := range reqs.Items {
		if err := req.Validate(); err != nil {
			apiErr := rest_err.NewBadRequestError(err.Error())
			logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
			return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
		}
	}

	updatedCount, apiErr := vc.service.BulkUpdateAltaiPhyItem(c.Context(), *claims, reqs.Items)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": updatedCount})
}

func (vc *altaiPhyCheckHandler) Finish(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	result, apiErr := vc.service.FinishCheck(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": result})
}

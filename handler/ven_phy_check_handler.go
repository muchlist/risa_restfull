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

func NewVenPhyCheckHandler(venPhyCheckService service.VenPhyCheckServiceAssumer) *venPhyCheckHandler {
	return &venPhyCheckHandler{
		service: venPhyCheckService,
	}
}

type venPhyCheckHandler struct {
	service service.VenPhyCheckServiceAssumer
}

func (vc *venPhyCheckHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := vc.service.InsertVenPhyCheck(c.Context(), *claims, req.Name, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": insertID})
}

func (vc *venPhyCheckHandler) InsertQuarter(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := vc.service.InsertVenPhyCheck(c.Context(), *claims, req.Name, true)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": insertID})
}

func (vc *venPhyCheckHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := vc.service.DeleteVenPhyCheck(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("cctv check %s berhasil dihapus", id)})
}

func (vc *venPhyCheckHandler) Get(c *fiber.Ctx) error {
	checkID := c.Params("id")

	check, apiErr := vc.service.GetVenPhyCheckByID(c.Context(), checkID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": check})
}

// Find menampilkan list check
// Query [branch, start, end, limit]
func (vc *venPhyCheckHandler) Find(c *fiber.Ctx) error {
	branch := c.Query("branch")
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	checkList, apiErr := vc.service.FindVenPhyCheck(c.Context(), branch, dto.FilterTimeRangeLimit{
		FilterStart: int64(start),
		FilterEnd:   int64(end),
		Limit:       int64(limit),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": checkList})
}

func (vc *venPhyCheckHandler) UpdateCheckItem(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.VenPhyCheckItemUpdateRequest
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

	checkUpdated, apiErr := vc.service.UpdateVenPhyCheckItem(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": checkUpdated})
}

func (vc *venPhyCheckHandler) BulkUpdateCheckItem(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var reqs dto.BulkVenPhyCheckUpdateRequest
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

	updatedCount, apiErr := vc.service.BulkUpdateVenPhyItem(c.Context(), *claims, reqs.Items)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": updatedCount})
}

func (vc *venPhyCheckHandler) Finish(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	result, apiErr := vc.service.FinishCheck(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": result})
}

func (vc *venPhyCheckHandler) FreshUpdateNameCCTV(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Params("branch")
	if branch == "" {
		branch = claims.Branch
	}

	result, apiErr := vc.service.FreshUpdateNameCCTV(c.Context(), branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": result})
}

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

func NewAltaiCheckHandler(altaiCheckService service.AltaiCheckServiceAssumer) *altaiCheckHandler {
	return &altaiCheckHandler{
		service: altaiCheckService,
	}
}

type altaiCheckHandler struct {
	service service.AltaiCheckServiceAssumer
}

func (ac *altaiCheckHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	insertID, apiErr := ac.service.InsertAltaiCheck(*claims)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": insertID})
}

func (ac *altaiCheckHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := ac.service.DeleteAltaiCheck(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("altai check %s berhasil dihapus", id)})
}

func (ac *altaiCheckHandler) Get(c *fiber.Ctx) error {
	checkID := c.Params("id")

	check, apiErr := ac.service.GetAltaiCheckByID(checkID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": check})
}

// Find menampilkan list check
// Query [branch, start, end, limit]
func (ac *altaiCheckHandler) Find(c *fiber.Ctx) error {
	branch := c.Query("branch")
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	checkList, apiErr := ac.service.FindAltaiCheck(branch, dto.FilterTimeRangeLimit{
		FilterStart: int64(start),
		FilterEnd:   int64(end),
		Limit:       int64(limit),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": checkList})
}

func (ac *altaiCheckHandler) UpdateCheckItem(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.AltaiCheckItemUpdateRequest
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

	checkUpdated, apiErr := ac.service.UpdateAltaiCheckItem(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": checkUpdated})
}

func (ac *altaiCheckHandler) BulkUpdateCheckItem(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var reqs dto.BulkAltaiCheckUpdateRequest
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

	updatedCount, apiErr := ac.service.BulkUpdateAltaiItem(*claims, reqs.Items)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": updatedCount})
}

func (ac *altaiCheckHandler) Finish(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	result, apiErr := ac.service.FinishCheck(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": result})
}

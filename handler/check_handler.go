package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"time"
)

func NewCheckHandler(checkService service.CheckServiceAssumer) *checkHandler {
	return &checkHandler{
		service: checkService,
	}
}

type checkHandler struct {
	service service.CheckServiceAssumer
}

func (ch *checkHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.CheckRequest
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

	insertID, apiErr := ch.service.InsertCheck(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan check berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// GetCheck menampilkan checkDetail
func (ch *checkHandler) GetCheck(c *fiber.Ctx) error {
	checkID := c.Params("id")

	check, apiErr := ch.service.GetCheckByID(c.Context(), checkID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": check})
}

// Find menampilkan list check
// Query [branch, start, end, limit]
func (ch *checkHandler) Find(c *fiber.Ctx) error {
	branch := c.Query("branch")
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	checkList, apiErr := ch.service.FindCheck(c.Context(), branch, dto.FilterTimeRangeLimit{
		FilterStart: int64(start),
		FilterEnd:   int64(end),
		Limit:       int64(limit),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": checkList})
}

func (ch *checkHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := ch.service.DeleteCheck(c.Context(), *claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("check %s berhasil dihapus", id)})
}

func (ch *checkHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	checkID := c.Params("id")

	var req dto.CheckEditRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	checkEdited, apiErr := ch.service.EditCheck(c.Context(), *claims, checkID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": checkEdited})
}

func (ch *checkHandler) UpdateCheckItem(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.CheckChildUpdateRequest
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

	checkUpdated, apiErr := ch.service.UpdateCheckItem(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": checkUpdated})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (ch *checkHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	childID := c.Params("child_id")

	// cek apakah ID check && branch ada
	_, apiErr := ch.service.GetCheckByID(c.Context(), id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%s%s%v", id, childID, time.Now().Unix())

	// simpan image
	pathInDB, apiErr := saveImage(c, *claims, "check", randomName, true)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	checkResult, apiErr := ch.service.PutChildImage(c.Context(), *claims, id, childID, pathInDB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": checkResult})
}

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

func NewCheckHandler(checkService service.CheckServiceAssumer) *checkHandler {
	return &checkHandler{
		service: checkService,
	}
}

type checkHandler struct {
	service service.CheckServiceAssumer
}

func (x *checkHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.CheckRequest
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

	insertID, apiErr := x.service.InsertCheck(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan check berhasil, ID: %s", *insertID)}
	return c.JSON(res)
}

// GetCheck menampilkan checkDetail
func (x *checkHandler) GetCheck(c *fiber.Ctx) error {

	checkID := c.Params("id")

	check, apiErr := x.service.GetCheckByID(checkID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(check)
}

// Find menampilkan list check
// Query [branch, start, end, limit]
func (x *checkHandler) Find(c *fiber.Ctx) error {

	branch := c.Query("branch")
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	checkList, apiErr := x.service.FindCheck(branch, dto.FilterTimeRangeLimit{
		FilterStart: int64(start),
		FilterEnd:   int64(end),
		Limit:       int64(limit),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"check_list": checkList})
}

func (x *checkHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := x.service.DeleteCheck(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": fmt.Sprintf("check %s berhasil dihapus", id)})
}

func (x *checkHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	checkID := c.Params("id")

	var req dto.CheckEditRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	checkEdited, apiErr := x.service.EditCheck(*claims, checkID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(checkEdited)
}

func (x *checkHandler) UpdateCheckItem(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.CheckChildUpdateRequest
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

	checkUpdated, apiErr := x.service.UpdateCheckItem(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(checkUpdated)
}

//UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (x *checkHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	childID := c.Params("child_id")

	// cek apakah ID check && branch ada
	_, apiErr := x.service.GetCheckByID(id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "check", id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	// update path image di database
	checkResult, apiErr := x.service.PutChildImage(*claims, id, childID, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(checkResult)
}

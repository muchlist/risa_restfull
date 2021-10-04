package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/statuses"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"time"
)

func NewOtherHandler(otherService service.OtherServiceAssumer) *otherHandler {
	return &otherHandler{
		service: otherService,
	}
}

type otherHandler struct {
	service service.OtherServiceAssumer
}

func (ot *otherHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.OtherRequest
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

	insertID, apiErr := ot.service.InsertOther(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan %s berhasil, ID: %s", req.SubCategory, *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// GetOther menampilkan otherDetail
func (ot *otherHandler) GetOther(c *fiber.Ctx) error {
	otherID := c.Params("id")

	other, apiErr := ot.service.GetOtherByID(c.Context(), otherID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": other})
}

// Find menampilkan list other
// Param [cat]
// Query [branch, name, ip, location, disable, division, seat]
func (ot *otherHandler) Find(c *fiber.Ctx) error {
	cat := c.Params("cat")
	branch := c.Query("branch")
	division := c.Query("division")
	name := c.Query("name")
	ip := c.Query("ip")
	location := c.Query("location")

	// cat validation
	if apiErr := subCategoryValidation(cat); apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	var disable bool
	if c.Query("disable") != "" {
		disable = true
	}

	if branch == "" {
		branch = c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim).Branch
	}

	filterA := dto.FilterOther{
		FilterBranch:      branch,
		FilterSubCategory: cat,
		FilterLocation:    location,
		FilterDivision:    division,
		FilterIP:          ip,
		FilterName:        name,
		FilterDisable:     disable,
	}

	otherList, generalList, apiErr := ot.service.FindOther(c.Context(), filterA)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fiber.Map{
		"other_list": otherList,
		"extra_list": generalList,
	}})
}

// DisableOther menghilangkan other dari list
// Param status [enable, disable] cat
func (ot *otherHandler) DisableOther(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	userID := c.Params("id")
	cat := c.Params("cat")
	status := c.Params("status")

	// cat validation
	if apiErr := subCategoryValidation(cat); apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// validation
	statusAvailable := []string{statuses.Disable, statuses.Enable}
	if !sfunc.InSlice(status, statusAvailable) {
		apiErr := rest_err.NewBadRequestError(fmt.Sprintf("Status yang dimasukkan tidak tersedia. gunakan %s", statusAvailable))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	var statusBool bool
	if status == statuses.Disable {
		statusBool = true
	}

	otherList, apiErr := ot.service.DisableOther(c.Context(), userID, *claims, cat, statusBool)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": otherList})
}

func (ot *otherHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	cat := c.Params("cat")

	// cat validation
	if apiErr := subCategoryValidation(cat); apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := ot.service.DeleteOther(c.Context(), *claims, cat, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("data %s berhasil dihapus", id)})
}

func (ot *otherHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	otherID := c.Params("id")

	var req dto.OtherEditRequest
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

	otherEdited, apiErr := ot.service.EditOther(c.Context(), *claims, otherID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": otherEdited})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (ot *otherHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	// cek apakah ID other && branch ada
	_, apiErr := ot.service.GetOtherByID(c.Context(), id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%s%v", id, time.Now().Unix())
	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "other", randomName, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	otherResult, apiErr := ot.service.PutImage(c.Context(), *claims, id, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": otherResult})
}

func subCategoryValidation(cat string) rest_err.APIError {
	if cat == "" {
		return rest_err.NewBadRequestError(fmt.Sprintf("param category wajib disertakan. gunakan %s", category.GetSubCategoryAvailable()))
	}

	if !sfunc.InSlice(cat, category.GetSubCategoryAvailable()) {
		return rest_err.NewBadRequestError(fmt.Sprintf("category yang dimasukkan tidak tersedia. gunakan %s", category.GetSubCategoryAvailable()))
	}
	return nil
}

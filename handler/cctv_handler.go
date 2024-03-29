package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/statuses"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/sfunc"
	"time"
)

func NewCctvHandler(cctvService service.CctvServiceAssumer) *cctvHandler {
	return &cctvHandler{
		service: cctvService,
	}
}

type cctvHandler struct {
	service service.CctvServiceAssumer
}

func (ctv *cctvHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.CctvRequest
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

	insertID, apiErr := ctv.service.InsertCctv(c.Context(), *claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan cctv berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// GetCctv menampilkan cctvDetail
func (ctv *cctvHandler) GetCctv(c *fiber.Ctx) error {
	cctvID := c.Params("id")

	cctv, apiErr := ctv.service.GetCctvByID(c.Context(), cctvID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": cctv})
}

// Find menampilkan list cctv
// Query [branch, name, ip, location, disable]
func (ctv *cctvHandler) Find(c *fiber.Ctx) error {
	branch := c.Query("branch")
	name := c.Query("name")
	ip := c.Query("ip")
	location := c.Query("location")
	var disable bool
	if c.Query("disable") != "" {
		disable = true
	}

	if branch == "" {
		branch = c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim).Branch
	}

	filterA := dto.FilterBranchLocIPNameDisable{
		FilterBranch:   branch,
		FilterLocation: location,
		FilterIP:       ip,
		FilterName:     name,
		FilterDisable:  disable,
	}

	cctvList, generalList, apiErr := ctv.service.FindCctv(c.Context(), filterA)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fiber.Map{
		"cctv_list":  cctvList,
		"extra_list": generalList,
	}})
}

// DisableCctv menghilangkan cctv dari list
// Param status [enable, disable]
func (ctv *cctvHandler) DisableCctv(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	userID := c.Params("id")
	status := c.Params("status")

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

	cctvList, apiErr := ctv.service.DisableCctv(c.Context(), userID, *claims, statusBool)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": cctvList})
}

func (ctv *cctvHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	forceStr := c.Query("force")
	var force bool
	if forceStr == "1" {
		force = true
	}

	apiErr := ctv.service.DeleteCctv(c.Context(), *claims, id, force)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("cctv %s berhasil dihapus", id)})
}

func (ctv *cctvHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	cctvID := c.Params("id")

	var req dto.CctvEditRequest
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

	cctvEdited, apiErr := ctv.service.EditCctv(c.Context(), *claims, cctvID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": cctvEdited})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (ctv *cctvHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	// cek apakah ID cctv && branch ada
	_, apiErr := ctv.service.GetCctvByID(c.Context(), id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%s%v", id, time.Now().Unix())
	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "cctv", randomName, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	cctvResult, apiErr := ctv.service.PutImage(c.Context(), *claims, id, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": cctvResult})
}

func (ctv *cctvHandler) Merge(c *fiber.Ctx) error {
	cctvID1 := c.Params("cctv1")
	cctvID2 := c.Params("cctv2")

	message, apiErr := ctv.service.MergeCctv(c.Context(), cctvID1, cctvID2)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": message})
}

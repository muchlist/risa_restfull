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

func NewComputerHandler(computerService service.ComputerServiceAssumer) *computerHandler {
	return &computerHandler{
		service: computerService,
	}
}

type computerHandler struct {
	service service.ComputerServiceAssumer
}

func (x *computerHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.ComputerRequest
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

	insertID, apiErr := x.service.InsertComputer(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan computer berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// GetComputer menampilkan computerDetail
func (x *computerHandler) GetComputer(c *fiber.Ctx) error {
	computerID := c.Params("id")

	computer, apiErr := x.service.GetComputerByID(computerID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": computer})
}

// Find menampilkan list computer
// Query [branch, name, ip, location, disable, division, seat]
func (x *computerHandler) Find(c *fiber.Ctx) error {
	branch := c.Query("branch")
	division := c.Query("division")
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

	seatManagement := -1
	if c.Query("seat") == "1" {
		seatManagement = 1
	}
	if c.Query("seat") == "0" {
		seatManagement = 0
	}

	filterA := dto.FilterComputer{
		FilterBranch:         branch,
		FilterLocation:       location,
		FilterDivision:       division,
		FilterIP:             ip,
		FilterName:           name,
		FilterDisable:        disable,
		FilterSeatManagement: seatManagement,
	}

	computerList, generalList, apiErr := x.service.FindComputer(filterA)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fiber.Map{
		"computer_list": computerList,
		"extra_list":    generalList,
	}})
}

// DisableComputer menghilangkan computer dari list
// Param status [enable, disable]
func (x *computerHandler) DisableComputer(c *fiber.Ctx) error {
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

	computerList, apiErr := x.service.DisableComputer(userID, *claims, statusBool)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": computerList})
}

func (x *computerHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := x.service.DeleteComputer(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("computer %s berhasil dihapus", id)})
}

func (x *computerHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	computerID := c.Params("id")

	var req dto.ComputerEditRequest
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

	computerEdited, apiErr := x.service.EditComputer(*claims, computerID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": computerEdited})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (x *computerHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	// cek apakah ID computer && branch ada
	_, apiErr := x.service.GetComputerByID(id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%s%v", id, time.Now().Unix())
	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "computer", randomName, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	computerResult, apiErr := x.service.PutImage(*claims, id, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": computerResult})
}

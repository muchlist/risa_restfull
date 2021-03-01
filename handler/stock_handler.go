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

func NewStockHandler(stockService service.StockServiceAssumer) *stockHandler {
	return &stockHandler{
		service: stockService,
	}
}

type stockHandler struct {
	service service.StockServiceAssumer
}

func (s *stockHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.StockRequest
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

	insertID, apiErr := s.service.InsertStock(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan stock berhasil, ID: %s", *insertID)}
	return c.JSON(res)
}

func (s *stockHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	stockID := c.Params("id")

	var req dto.StockEditRequest
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

	stockEdited, apiErr := s.service.EditStock(*claims, stockID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(stockEdited)
}

func (s *stockHandler) ChangeQty(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	stockID := c.Params("id")

	var req dto.StockChangeRequest
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

	stockEdited, apiErr := s.service.ChangeQtyStock(*claims, stockID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(stockEdited)
}

// GetStock menampilkan stock Detail
func (s *stockHandler) GetStock(c *fiber.Ctx) error {

	stockID := c.Params("id")

	stock, apiErr := s.service.GetStockByID(stockID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(stock)
}

// Find menampilkan list stock
// Query [branch, name, category, disable]
func (s *stockHandler) Find(c *fiber.Ctx) error {

	branch := c.Query("branch")
	name := c.Query("name")
	category := c.Query("category")
	var disable bool
	if c.Query("disable") != "" {
		disable = true
	}

	filterA := dto.FilterBranchNameCatDisable{
		FilterBranch:   branch,
		FilterName:     name,
		FilterCategory: category,
		FilterDisable:  disable,
	}

	stockList, apiErr := s.service.FindStock(filterA)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"stock_list": stockList})
}

//DisableStock menghilangkan stock dari list
// Param status [enable, disable]
func (s *stockHandler) DisableStock(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	userID := c.Params("id")
	status := c.Params("status")

	// validation
	statusAvailable := []string{"disable", "enable"}
	if !sfunc.InSlice(status, statusAvailable) {
		apiErr := rest_err.NewBadRequestError(fmt.Sprintf("Status yang dimasukkan tidak tersedia. gunakan %s", statusAvailable))
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	var statusBool bool
	if status == "disable" {
		statusBool = true
	}

	stockList, apiErr := s.service.DisableStock(userID, *claims, statusBool)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"stock_list": stockList})
}

func (s *stockHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := s.service.DeleteStock(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": fmt.Sprintf("stock %s berhasil dihapus", id)})
}

//UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (s *stockHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	// cek apakah ID stock && branch ada
	_, apiErr := s.service.GetStockByID(id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "stock", id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	// update path image di database
	stockResult, apiErr := s.service.PutImage(*claims, id, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(stockResult)
}

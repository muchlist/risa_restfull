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

func NewCheckItemHandler(checkItemService service.CheckItemServiceAssumer) *checkItemHandler {
	return &checkItemHandler{
		service: checkItemService,
	}
}

type checkItemHandler struct {
	service service.CheckItemServiceAssumer
}

func (i *checkItemHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.CheckItemRequest
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

	insertID, apiErr := i.service.InsertCheckItem(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan checkItem berhasil, ID: %s", *insertID)}
	return c.JSON(res)
}

// GetCheckItem menampilkan checkItemDetail
func (i *checkItemHandler) GetCheckItem(c *fiber.Ctx) error {

	checkItemID := c.Params("id")

	checkItem, apiErr := i.service.GetCheckItemByID(checkItemID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(checkItem)
}

// Find menampilkan list checkItem
// Query [branch, name, problem, disable]
func (i *checkItemHandler) Find(c *fiber.Ctx) error {

	branch := c.Query("branch")
	name := c.Query("name")
	var disable bool
	if c.Query("disable") != "" {
		disable = true
	}
	var haveProblem bool
	if c.Query("problem") != "" {
		disable = true
	}

	filterA := dto.FilterBranchNameDisable{
		FilterBranch:  branch,
		FilterName:    name,
		FilterDisable: disable,
	}

	checkItemList, apiErr := i.service.FindCheckItem(filterA, haveProblem)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"check_item_list": checkItemList})
}

//DisableCheckItem menghilangkan checkItem dari list
// Param status [enable, disable]
func (i *checkItemHandler) DisableCheckItem(c *fiber.Ctx) error {
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

	checkItemList, apiErr := i.service.DisableCheckItem(userID, *claims, statusBool)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"checkItem_list": checkItemList})
}

func (i *checkItemHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := i.service.DeleteCheckItem(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": fmt.Sprintf("checkItem %s berhasil dihapus", id)})
}

func (i *checkItemHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	checkItemID := c.Params("id")

	var req dto.CheckItemEditRequest
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

	checkItemEdited, apiErr := i.service.EditCheckItem(*claims, checkItemID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(checkItemEdited)
}

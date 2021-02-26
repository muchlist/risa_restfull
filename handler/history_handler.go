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

func NewHistoryHandler(histService service.HistoryServiceAssumer) *historyHandler {
	return &historyHandler{
		service: histService,
	}
}

type historyHandler struct {
	service service.HistoryServiceAssumer
}

func (h *historyHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.HistoryRequest
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

	insertID, apiErr := h.service.InsertHistory(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan history berhasi, ID: %s", *insertID)}
	return c.JSON(res)
}

func (h *historyHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	historyID := c.Params("id")

	var req dto.HistoryEditRequest
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

	historyEdited, apiErr := h.service.EditHistory(*claims, historyID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}
	return c.JSON(historyEdited)
}

// Find menampilkan list history
// Query [branch, category, c_status, start, end, limit]
func (h *historyHandler) Find(c *fiber.Ctx) error {

	branch := c.Query("branch")
	category := c.Query("category")
	cStatus := stringToInt(c.Query("c_status"))
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	filterA := dto.FilterBranchCatComplete{
		Branch:         branch,
		Category:       category,
		CompleteStatus: cStatus,
	}

	filterB := dto.FilterTimeRangeLimit{
		Start: int64(start),
		End:   int64(end),
		Limit: int64(limit),
	}

	histories, apiErr := h.service.FindHistory(filterA, filterB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"histories": histories})
}

// Find menampilkan list history berdasarkan parent string
func (h *historyHandler) FindFromParent(c *fiber.Ctx) error {

	parentID := c.Params("id")

	histories, apiErr := h.service.FindHistoryForParent(parentID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"histories": histories})
}

// Find menampilkan list history berdasarkan user
// Query [start, end, limit]
func (h *historyHandler) FindFromUser(c *fiber.Ctx) error {

	userID := c.Params("id")

	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	filter := dto.FilterTimeRangeLimit{
		Start: int64(start),
		End:   int64(end),
		Limit: int64(limit),
	}

	histories, apiErr := h.service.FindHistoryForUser(userID, filter)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"histories": histories})
}

// GetHistory menampilkan historyDetail
func (h *historyHandler) GetHistory(c *fiber.Ctx) error {

	userID := c.Params("id")

	history, apiErr := h.service.GetHistory(userID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(history)
}

func (h *historyHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := h.service.DeleteHistory(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"msg": fmt.Sprintf("history %s berhasil dihapus", id)})
}

//UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (h *historyHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	// cek apakah ID history && branch ada
	_, apiErr := h.service.GetHistory(id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "history", id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	// update path image di database
	cctvResult, apiErr := h.service.PutImage(*claims, id, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(cctvResult)
}

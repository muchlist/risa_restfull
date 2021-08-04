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
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := h.service.InsertHistory(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan history berhasi, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (h *historyHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	historyID := c.Params("id")

	var req dto.HistoryEditRequest
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

	historyEdited, apiErr := h.service.EditHistory(*claims, historyID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": historyEdited})
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
		FilterBranch:         branch,
		FilterCategory:       category,
		FilterCompleteStatus: cStatus,
	}

	filterB := dto.FilterTimeRangeLimit{
		FilterStart: int64(start),
		FilterEnd:   int64(end),
		Limit:       int64(limit),
	}

	histories, apiErr := h.service.FindHistory(filterA, filterB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": histories})
}

// FindUnwind menampilkan list history unwind
// Query [branch, category, c_status, start, end, limit]
func (h *historyHandler) FindUnwind(c *fiber.Ctx) error {
	branch := c.Query("branch")
	category := c.Query("category")
	cStatus := stringToInt(c.Query("c_status"))
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	filterA := dto.FilterBranchCatComplete{
		FilterBranch:         branch,
		FilterCategory:       category,
		FilterCompleteStatus: cStatus,
	}

	filterB := dto.FilterTimeRangeLimit{
		FilterStart: int64(start),
		FilterEnd:   int64(end),
		Limit:       int64(limit),
	}

	histories, apiErr := h.service.UnwindHistory(filterA, filterB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": histories})
}

// Find menampilkan list history berdasarkan parent string
func (h *historyHandler) FindFromParent(c *fiber.Ctx) error {
	parentID := c.Params("id")

	histories, apiErr := h.service.FindHistoryForParent(parentID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": histories})
}

// Find menampilkan list history berdasarkan user
// Query [start, end, limit]
func (h *historyHandler) FindFromUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))
	limit := stringToInt(c.Query("limit"))

	filter := dto.FilterTimeRangeLimit{
		FilterStart: int64(start),
		FilterEnd:   int64(end),
		Limit:       int64(limit),
	}

	histories, apiErr := h.service.FindHistoryForUser(userID, filter)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": histories})
}

// GetHistory menampilkan historyDetail
func (h *historyHandler) GetHistory(c *fiber.Ctx) error {
	userID := c.Params("id")

	history, apiErr := h.service.GetHistory(userID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": history})
}

func (h *historyHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := h.service.DeleteHistory(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("history %s berhasil dihapus", id)})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (h *historyHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	// cek apakah ID history && branch ada
	_, apiErr := h.service.GetHistory(id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%s%v", id, time.Now().Unix())
	// simpan image
	pathInDB, apiErr := saveImage(c, *claims, "history", randomName, true)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	historyResult, apiErr := h.service.PutImage(*claims, id, pathInDB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": historyResult})
}

// UploadImageWithoutParent melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan mengembalikan nama image
func (h *historyHandler) UploadImageWithoutParent(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	randomName := fmt.Sprintf("%s%v", claims.Identity, time.Now().Unix())
	// simpan image
	pathInDB, apiErr := saveImage(c, *claims, "history", randomName, true)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": pathInDB})
}

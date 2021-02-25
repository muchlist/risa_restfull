package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"strconv"
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
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	if err := req.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	insertID, apiErr := h.service.InsertHistory(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	res := fiber.Map{"msg": fmt.Sprintf("Menambahkan history berhasi, ID: %s", *insertID)}
	return c.JSON(res)
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

	history, apiErr := h.service.GetHistory(userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(history)
}

func stringToInt(queryString string) int {
	number, err := strconv.Atoi(queryString)
	if err != nil {
		return 0
	}
	return number
}

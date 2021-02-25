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

//Find menampilkan list user
func (h *historyHandler) Find(c *fiber.Ctx) error {

	branch := c.Query("branch")
	category := c.Query("category")
	cStatus := queryToInt(c, "complete_status")
	start := queryToInt(c, "start")
	end := queryToInt(c, "end")
	limit := queryToInt(c, "limit")

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

func queryToInt(c *fiber.Ctx, queryKey string) int {
	inted, err := strconv.Atoi(c.Query(queryKey))
	if err != nil {
		return 0
	}
	return inted
}

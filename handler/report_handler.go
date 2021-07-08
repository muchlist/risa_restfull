package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/risa_restfull/service"
)

func NewReportHandler(histService service.ReportServiceAssumer) *reportHandler {
	return &reportHandler{
		service: histService,
	}
}

type reportHandler struct {
	service service.ReportServiceAssumer
}

// GeneratePDF membuat pdf
// Query [branch, start, end]
func (h *reportHandler) GeneratePDF(c *fiber.Ctx) error {
	branch := c.Query("branch")
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))

	fileName, apiErr := h.service.GenerateReportPDF(branch, int64(start), int64(end))
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": fileName})
}

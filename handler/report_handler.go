package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/mjwt"
	"github.com/muchlist/risa_restfull/utils/timegen"
	"net/http"
	"time"
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
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Query("branch")
	if branch == "" {
		branch = claims.Branch
	}

	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))

	pdfName, err2 := timegen.GetTimeAsName(int64(end))
	if err2 != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": rest_err.NewBadRequestError("gagal membuat nama pdf"), "data": nil})
	}
	pdfName = fmt.Sprintf("manual-%s", pdfName)

	_, apiErr := h.service.GenerateReportPDF(pdfName, branch, int64(start), int64(end))
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(dto.PdfFile{
		CreatedAt: time.Now().Unix(),
		CreatedBy: claims.Name,
		Branch:    branch,
		Name:      pdfName,
		Type:      "LAPORAN",
		FileName:  fmt.Sprintf("pdf/%s.pdf", pdfName),
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf/%s.pdf", pdfName)})
}

func (h *reportHandler) FindPDF(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Query("branch")
	if branch == "" {
		branch = claims.Branch
	}

	pdfList, apiErr := h.service.FindPdf(branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	if pdfList == nil {
		pdfList = []dto.PdfFile{}
	}
	return c.JSON(fiber.Map{"error": nil, "data": pdfList})
}

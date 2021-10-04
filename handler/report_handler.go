package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/pdftype"
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
	pdfName = fmt.Sprintf("support-%s", pdfName)

	_, apiErr := h.service.GenerateReportPDF(c.Context(), pdfName, branch, int64(start), int64(end))
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(c.Context(), dto.PdfFile{
		CreatedAt:     time.Now().Unix(),
		CreatedBy:     claims.Name,
		Branch:        branch,
		Name:          pdfName,
		Type:          "LAPORAN",
		FileName:      fmt.Sprintf("pdf/%s.pdf", pdfName),
		EndReportTime: int64(end),
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf/%s.pdf", pdfName)})
}

// GeneratePDFStartFromLast membuat pdf berdasarkan tanggal pdf sebelumnya dijadikan awal
// dan tanggal saat ini dijadikan akhir
// Query [branch]
func (h *reportHandler) GeneratePDFStartFromLast(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Query("branch")
	if branch == "" {
		branch = claims.Branch
	}
	currentTime := time.Now().Unix()

	pdfName, err2 := timegen.GetTimeAsName(currentTime)
	if err2 != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": rest_err.NewBadRequestError("gagal membuat nama pdf"), "data": nil})
	}
	pdfName = fmt.Sprintf("support-%s", pdfName)

	_, apiErr := h.service.GenerateReportPDFStartFromLast(c.Context(), pdfName, branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(c.Context(), dto.PdfFile{
		CreatedAt:     currentTime,
		CreatedBy:     claims.Name,
		Branch:        branch,
		Name:          pdfName,
		Type:          pdftype.Laporan,
		FileName:      fmt.Sprintf("pdf/%s.pdf", pdfName),
		EndReportTime: currentTime,
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf/%s.pdf", pdfName)})
}

// GeneratePDFVendor membuat pdf
// Query [branch, start, end]
func (h *reportHandler) GeneratePDFVendor(c *fiber.Ctx) error {
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
	pdfName = fmt.Sprintf("vendor-%s", pdfName)

	_, apiErr := h.service.GenerateReportPDFVendor(c.Context(), pdfName, branch, int64(start), int64(end))
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(c.Context(), dto.PdfFile{
		CreatedAt:     time.Now().Unix(),
		CreatedBy:     claims.Name,
		Branch:        branch,
		Name:          pdfName,
		Type:          pdftype.VendorSum,
		FileName:      fmt.Sprintf("pdf-vendor/%s.pdf", pdfName),
		EndReportTime: int64(end),
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf-vendor/%s.pdf", pdfName)})
}

// GeneratePDFVendorStartFromLast membuat pdf berdasarkan tanggal pdf sebelumnya dijadikan awal
// dan tanggal saat ini dijadikan akhir
// Query [branch]
func (h *reportHandler) GeneratePDFVendorStartFromLast(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Query("branch")
	if branch == "" {
		branch = claims.Branch
	}

	currentTime := time.Now().Unix()

	pdfName, err2 := timegen.GetTimeAsName(currentTime)
	if err2 != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": rest_err.NewBadRequestError("gagal membuat nama pdf"), "data": nil})
	}
	pdfName = fmt.Sprintf("vendor-%s", pdfName)

	_, apiErr := h.service.GenerateReportPDFVendorStartFromLast(c.Context(), pdfName, branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(c.Context(), dto.PdfFile{
		CreatedAt:     currentTime,
		CreatedBy:     claims.Name,
		Branch:        branch,
		Name:          pdfName,
		Type:          pdftype.VendorSum,
		FileName:      fmt.Sprintf("pdf-vendor/%s.pdf", pdfName),
		EndReportTime: currentTime,
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf-vendor/%s.pdf", pdfName)})
}

func (h *reportHandler) GeneratePDFDailyReportVendor(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Query("branch")
	if branch == "" {
		branch = claims.Branch
	}

	start := int64(stringToInt(c.Query("start")))
	end := int64(stringToInt(c.Query("end")))

	target := int64(stringToInt(c.Query("target")))
	if target != 0 {
		end = target
	}

	currentTime := time.Now().Unix()
	if end > currentTime {
		end = currentTime
	}

	fmt.Printf("start : %d, end : %d\n", start, end)

	if end == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": rest_err.NewBadRequestError("target waktu pdf harus ditentukan (target=12345678)"), "data": nil})
	}

	pdfName, err2 := timegen.GetTimeAsName(end)
	if err2 != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": rest_err.NewBadRequestError("gagal membuat nama pdf"), "data": nil})
	}
	pdfName = fmt.Sprintf("daily-vendor-%s", pdfName)

	_, apiErr := h.service.GenerateReportVendorDaily(c.Context(), pdfName, branch, start, end)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(c.Context(), dto.PdfFile{
		CreatedAt:     currentTime,
		CreatedBy:     claims.Name,
		Branch:        branch,
		Name:          pdfName,
		Type:          pdftype.Vendor,
		FileName:      fmt.Sprintf("pdf-vendor/%s.pdf", pdfName),
		EndReportTime: currentTime,
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf-vendor/%s.pdf", pdfName)})
}

// GeneratePDFVendorDailyStartFromLast membuat pdf berdasarkan tanggal pdf sebelumnya dijadikan awal
// dan tanggal saat ini dijadikan akhir
// Query [branch]
func (h *reportHandler) GeneratePDFVendorDailyStartFromLast(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Query("branch")
	if branch == "" {
		branch = claims.Branch
	}

	currentTime := time.Now().Unix()

	pdfName, err2 := timegen.GetTimeAsName(currentTime)
	if err2 != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": rest_err.NewBadRequestError("gagal membuat nama pdf"), "data": nil})
	}
	pdfName = fmt.Sprintf("daily-vendor-%s", pdfName)

	_, apiErr := h.service.GenerateReportVendorDailyStartFromLast(c.Context(), pdfName, branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(c.Context(), dto.PdfFile{
		CreatedAt:     currentTime,
		CreatedBy:     claims.Name,
		Branch:        branch,
		Name:          pdfName,
		Type:          pdftype.Vendor,
		FileName:      fmt.Sprintf("pdf-vendor/%s.pdf", pdfName),
		EndReportTime: currentTime,
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf-vendor/%s.pdf", pdfName)})
}

// GeneratePDFVendorMonthly membuat pdf
// Query [branch, start, end]
func (h *reportHandler) GeneratePDFVendorMonthly(c *fiber.Ctx) error {
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
	pdfName = fmt.Sprintf("vendor-monthly%s", pdfName)

	_, apiErr := h.service.GenerateReportPDFVendorMonthly(c.Context(), pdfName, branch, int64(start), int64(end))
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(c.Context(), dto.PdfFile{
		CreatedAt:     time.Now().Unix(),
		CreatedBy:     claims.Name,
		Branch:        branch,
		Name:          pdfName,
		Type:          pdftype.VendorMonthly,
		FileName:      fmt.Sprintf("pdf-v-month/%s.pdf", pdfName),
		EndReportTime: int64(end),
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf-v-month/%s.pdf", pdfName)})
}

// GeneratePDFStock membuat pdf
// Query [branch, start, end]
func (h *reportHandler) GeneratePDFStock(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Query("branch")
	if branch == "" {
		branch = claims.Branch
	}

	category := c.Query("category")
	start := stringToInt(c.Query("start"))
	end := stringToInt(c.Query("end"))

	pdfName, err2 := timegen.GetTimeAsName(int64(end))
	if err2 != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": rest_err.NewBadRequestError("gagal membuat nama pdf"), "data": nil})
	}
	pdfName = fmt.Sprintf("stock-%s", pdfName)

	_, apiErr := h.service.GenerateStockReportRestock(c.Context(), pdfName, branch, category, int64(start), int64(end))
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// simpan pdf ke database
	_, apiErr = h.service.InsertPdf(c.Context(), dto.PdfFile{
		CreatedAt:     time.Now().Unix(),
		CreatedBy:     claims.Name,
		Branch:        branch,
		Name:          pdfName,
		Type:          pdftype.Stock,
		FileName:      fmt.Sprintf("pdf-stock/%s.pdf", pdfName),
		EndReportTime: int64(end),
	})

	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("pdf-stock/%s.pdf", pdfName)})
}

func (h *reportHandler) FindPDF(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	branch := c.Query("branch")
	pdfType := c.Query("type")
	if branch == "" {
		branch = claims.Branch
	}

	pdfList, apiErr := h.service.FindPdf(c.Context(), branch, pdfType)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	if pdfList == nil {
		pdfList = []dto.PdfFile{}
	}
	return c.JSON(fiber.Map{"error": nil, "data": pdfList})
}

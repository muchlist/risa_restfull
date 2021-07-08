package service

import (
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/reportdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/pdfgen"
)

func NewReportService(
	histDao historydao.HistoryDaoAssumer,
	checkDao checkdao.CheckDaoAssumer,
	pdfDao reportdao.PdfDaoAssumer,
) ReportServiceAssumer {
	return &reportService{
		daoH: histDao,
		daoC: checkDao,
		daoP: pdfDao,
	}
}

type reportService struct {
	daoH historydao.HistoryDaoAssumer
	daoC checkdao.CheckDaoAssumer
	daoP reportdao.PdfDaoAssumer
}

type ReportServiceAssumer interface {
	InsertPdf(input dto.PdfFile) (*string, rest_err.APIError)
	GenerateReportPDF(name string, branch string, start int64, end int64) (*string, rest_err.APIError)
	FindPdf(branch string) ([]dto.PdfFile, rest_err.APIError)
}

// GeneratePDFReport
func (r *reportService) GenerateReportPDF(name string, branch string, start int64, end int64) (*string, rest_err.APIError) {
	if start > end && start < 0 {
		return nil, rest_err.NewBadRequestError("tanggal terakhir tidak boleh lebih besar dari tanggal awal")
	}

	// GET HISTORIES
	historyList, err := r.daoH.FindHistoryForReport(branch, start, end)
	if err != nil {
		return nil, err
	}

	// GET CHECK LIST
	checkList, err := r.daoC.FindCheckForReports(branch, dto.FilterTimeRangeLimit{
		FilterStart: start,
		FilterEnd:   end,
		Limit:       2,
	})
	if err != nil {
		return nil, err
	}

	errPDF := pdfgen.GeneratePDF(pdfgen.PDFReq{
		Name:        name,
		HistoryList: historyList,
		CheckList:   checkList,
		Start:       start,
		End:         end,
	})
	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat pdf", errPDF)
	}

	return &name, nil
}

func (r *reportService) InsertPdf(input dto.PdfFile) (*string, rest_err.APIError) {
	return r.daoP.InsertPdf(input)
}

func (r *reportService) FindPdf(branch string) ([]dto.PdfFile, rest_err.APIError) {
	return r.daoP.FindPdf(branch)
}

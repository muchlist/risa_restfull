package service

import (
	"fmt"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/pdfgen"
	"github.com/muchlist/risa_restfull/utils/timegen"
)

func NewReportService(
	histDao historydao.HistoryDaoAssumer,
	checkDao checkdao.CheckDaoAssumer,
) ReportServiceAssumer {
	return &reportService{
		daoH: histDao,
		daoC: checkDao,
	}
}

type reportService struct {
	daoH historydao.HistoryDaoAssumer
	daoC checkdao.CheckDaoAssumer
}
type ReportServiceAssumer interface {
	GenerateReportPDF(branch string, start int64, end int64) (*string, rest_err.APIError)
}

// GeneratePDFReport
func (r *reportService) GenerateReportPDF(branch string, start int64, end int64) (*string, rest_err.APIError) {
	if start > end && start < 0 {
		return nil, rest_err.NewBadRequestError("tanggal terakhir tidak boleh lebih besar dari tanggal awal")
	}

	pdfName, err2 := timegen.GetTimeAsName(end)
	if err2 != nil {
		return nil, rest_err.NewBadRequestError("gagal membuat nama pdf berdasarkan tanggal terakhir")
	}
	pdfName = fmt.Sprintf("manual-%s", pdfName)

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
		Name:        pdfName,
		HistoryList: historyList,
		CheckList:   checkList,
		Start:       start,
		End:         end,
	})
	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat pdf", errPDF)
	}

	return &pdfName, nil
}

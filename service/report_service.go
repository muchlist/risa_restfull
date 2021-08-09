package service

import (
	"fmt"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/reportdao"
	"github.com/muchlist/risa_restfull/dao/vendorcheckdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/pdfgen"
)

func NewReportService(
	histDao historydao.HistoryDaoAssumer,
	checkDao checkdao.CheckDaoAssumer,
	checkVendorDao vendorcheckdao.CheckVendorDaoAssumer,
	pdfDao reportdao.PdfDaoAssumer,
) ReportServiceAssumer {
	return &reportService{
		daoH:  histDao,
		daoC:  checkDao,
		daoCV: checkVendorDao,
		daoP:  pdfDao,
	}
}

type reportService struct {
	daoH  historydao.HistoryDaoAssumer
	daoC  checkdao.CheckDaoAssumer
	daoCV vendorcheckdao.CheckVendorDaoAssumer
	daoP  reportdao.PdfDaoAssumer
}

type ReportServiceAssumer interface {
	InsertPdf(input dto.PdfFile) (*string, rest_err.APIError)
	GenerateReportPDF(name string, branch string, start int64, end int64) (*string, rest_err.APIError)
	GenerateReportPDFVendor(name string, branch string, start int64, end int64) (*string, rest_err.APIError)
	FindPdf(branch string) ([]dto.PdfFile, rest_err.APIError)
}

// GenerateReportPDF membuat laporan untuk it support
func (r *reportService) GenerateReportPDF(name string, branch string, start int64, end int64) (*string, rest_err.APIError) {
	if start > end && start < 0 {
		return nil, rest_err.NewBadRequestError("tanggal awal tidak boleh lebih besar dari tanggal akhir")
	}

	// GET HISTORIES 0, 4 sesuai start end inputan
	historyList04, err := r.daoH.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       "",
			FilterCompleteStatus: "0,4",
		}, dto.FilterTimeRangeLimit{
			FilterStart: start,
			FilterEnd:   end,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.daoH.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       "",
			FilterCompleteStatus: "1,2,3",
		}, dto.FilterTimeRangeLimit{
			FilterStart: end - (3 * 30 * 24 * 60 * 60), // 3 bulan,
			FilterEnd:   end,
			Limit:       300,
		},
	)

	historiesCombined := append(historyList04, historyList123...)

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
		HistoryList: historiesCombined,
		CheckList:   checkList,
		Start:       start,
		End:         end,
	})
	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat pdf", errPDF)
	}

	return &name, nil
}

// GenerateReportPDFVendor membuat pdf untuk vendor multinet
func (r *reportService) GenerateReportPDFVendor(name string, branch string, start int64, end int64) (*string, rest_err.APIError) {
	if start > end && start < 0 {
		return nil, rest_err.NewBadRequestError("tanggal awal tidak boleh lebih besar dari tanggal akhir")
	}

	// GET HISTORIES 0, 4 sesuai start end inputan
	historyList04, err := r.daoH.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: "0,4",
		}, dto.FilterTimeRangeLimit{
			FilterStart: start,
			FilterEnd:   end,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.daoH.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: "1,2,3",
		}, dto.FilterTimeRangeLimit{
			FilterStart: end - (3 * 30 * 24 * 60 * 60), // 3 bulan,
			FilterEnd:   end,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	historiesCombined := append(historyList04, historyList123...)

	vendorCheckList, err := r.daoCV.FindCheck(branch, dto.FilterTimeRangeLimit{
		FilterStart: start - (2 * 30 * 24 * 60 * 60), // batas awalnya di kurangi 2 bulan
		FilterEnd:   end,
		Limit:       20,
	}, true)
	if err != nil {
		return nil, err
	}

	errPDF := pdfgen.GeneratePDFVendor(pdfgen.PDFVendorReq{
		Name:            name,
		HistoryList:     historiesCombined,
		VendorCheckList: vendorCheckList,
		Start:           start,
		End:             end,
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

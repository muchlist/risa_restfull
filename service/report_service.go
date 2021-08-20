package service

import (
	"fmt"
	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/pdftype"
	"github.com/muchlist/risa_restfull/dao/altaicheckdao"
	"github.com/muchlist/risa_restfull/dao/altaiphycheckdao"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/reportdao"
	"github.com/muchlist/risa_restfull/dao/vendorcheckdao"
	"github.com/muchlist/risa_restfull/dao/venphycheckdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/pdfgen"
	"time"
)

// ReportParams berisi semua dao yang diperlukan reports service, karena sangat banyak maka dibuat struct
type ReportParams struct {
	History       historydao.HistoryDaoAssumer
	CheckIT       checkdao.CheckDaoAssumer
	CheckCCTV     vendorcheckdao.CheckVendorDaoAssumer
	CheckCCTVPhy  venphycheckdao.CheckVenPhyDaoAssumer
	CheckAltai    altaicheckdao.CheckAltaiDaoAssumer
	CheckAltaiPhy altaiphycheckdao.CheckAltaiPhyDaoAssumer
	Pdf           reportdao.PdfDaoAssumer
}

func NewReportService(dao ReportParams) ReportServiceAssumer {
	return &reportService{dao: dao}
}

type reportService struct {
	dao ReportParams
}

type ReportServiceAssumer interface {
	InsertPdf(input dto.PdfFile) (*string, rest_err.APIError)
	GenerateReportPDF(name string, branch string, start int64, end int64) (*string, rest_err.APIError)
	GenerateReportPDFStartFromLast(name string, branch string) (*string, rest_err.APIError)
	GenerateReportPDFVendor(name string, branch string, start int64, end int64) (*string, rest_err.APIError)
	GenerateReportPDFVendorStartFromLast(name string, branch string) (*string, rest_err.APIError)
	FindPdf(branch string, typePdf string) ([]dto.PdfFile, rest_err.APIError)
	GenerateReportVendorDaily(name string, branch string, target int64) (*string, rest_err.APIError)
}

// GenerateReportPDF membuat laporan untuk it support
func (r *reportService) GenerateReportPDF(name string, branch string, start int64, end int64) (*string, rest_err.APIError) {
	if start > end && start < 0 {
		return nil, rest_err.NewBadRequestError("tanggal awal tidak boleh lebih besar dari tanggal akhir")
	}

	currentTime := time.Now().Unix()
	if end > currentTime {
		end = currentTime
	}

	// GET HISTORIES 0, 4 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
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
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       "",
			FilterCompleteStatus: "1,2,3",
		}, dto.FilterTimeRangeLimit{
			FilterStart: start - (3 * 30 * 24 * 60 * 60), // 3 bulan,
			FilterEnd:   end,
			Limit:       300,
		},
	)

	historiesCombined := append(historyList04, historyList123...)

	// GET CHECK LIST
	checkList, err := r.dao.CheckIT.FindCheckForReports(branch, dto.FilterTimeRangeLimit{
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
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

// GenerateReportPDF membuat laporan untuk it support
func (r *reportService) GenerateReportPDFStartFromLast(name string, branch string) (*string, rest_err.APIError) {

	currentTime := time.Now().Unix()
	lastPDF, err := r.dao.Pdf.FindLastPdf(branch, pdftype.Laporan)
	if err != nil {
		return nil, rest_err.NewBadRequestError("Gagal mendapatkan data laporan sebelumnya")
	}
	lastPDFEndTime := lastPDF.EndReportTime

	if currentTime-lastPDFEndTime < 60*2 {
		return nil, rest_err.NewBadRequestError("Gagal. Jarak pembuatan laporan tidak boleh kurang dari 2 menit!")
	}

	// GET HISTORIES 0, 4 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       "",
			FilterCompleteStatus: "0,4",
		}, dto.FilterTimeRangeLimit{
			FilterStart: lastPDFEndTime,
			FilterEnd:   currentTime,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       "",
			FilterCompleteStatus: "1,2,3",
		}, dto.FilterTimeRangeLimit{
			FilterStart: lastPDFEndTime - (3 * 30 * 24 * 60 * 60), // 3 bulan,
			FilterEnd:   currentTime,
			Limit:       300,
		},
	)

	historiesCombined := append(historyList04, historyList123...)

	// GET CHECK LIST
	checkList, err := r.dao.CheckIT.FindCheckForReports(branch, dto.FilterTimeRangeLimit{
		FilterStart: lastPDFEndTime,
		FilterEnd:   currentTime,
	})
	if err != nil {
		return nil, err
	}

	errPDF := pdfgen.GeneratePDF(pdfgen.PDFReq{
		Name:        name,
		HistoryList: historiesCombined,
		CheckList:   checkList,
		Start:       lastPDFEndTime,
		End:         currentTime,
	})
	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

// GenerateReportPDFVendor membuat Pdf untuk vendor multinet
func (r *reportService) GenerateReportPDFVendor(name string, branch string, start int64, end int64) (*string, rest_err.APIError) {
	if start > end && start < 0 {
		return nil, rest_err.NewBadRequestError("tanggal awal tidak boleh lebih besar dari tanggal akhir")
	}

	currentTime := time.Now().Unix()
	if end > currentTime {
		end = currentTime
	}

	// GET HISTORIES 0, 4 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
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
	historyList123, err := r.dao.History.UnwindHistory(
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

	vendorCheckList, err := r.dao.CheckCCTV.FindCheck(branch, dto.FilterTimeRangeLimit{
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
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

// GenerateReportPDFVendor membuat Pdf untuk vendor multinet
func (r *reportService) GenerateReportPDFVendorStartFromLast(name string, branch string) (*string, rest_err.APIError) {

	currentTime := time.Now().Unix()
	lastPDF, err := r.dao.Pdf.FindLastPdf(branch, pdftype.Vendor)
	if err != nil {
		return nil, rest_err.NewBadRequestError("gagal mendapatkan data laporan sebelumnya")
	}
	lastPDFEndTime := lastPDF.EndReportTime

	if currentTime-lastPDFEndTime < 60*2 {
		return nil, rest_err.NewBadRequestError("Gagal. Jarak pembuatan laporan tidak boleh kurang dari 2 menit!")
	}

	// GET HISTORIES 0, 4 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: "0,4",
		}, dto.FilterTimeRangeLimit{
			FilterStart: lastPDFEndTime,
			FilterEnd:   currentTime,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: "1,2,3",
		}, dto.FilterTimeRangeLimit{
			FilterStart: lastPDFEndTime - (3 * 30 * 24 * 60 * 60), // 3 bulan,
			FilterEnd:   currentTime,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	historiesCombined := append(historyList04, historyList123...)

	vendorCheckList, err := r.dao.CheckCCTV.FindCheck(branch, dto.FilterTimeRangeLimit{
		FilterStart: lastPDFEndTime - (2 * 30 * 24 * 60 * 60), // batas awalnya di kurangi 2 bulan
		FilterEnd:   currentTime,
		Limit:       20,
	}, true)
	if err != nil {
		return nil, err
	}

	errPDF := pdfgen.GeneratePDFVendor(pdfgen.PDFVendorReq{
		Name:            name,
		HistoryList:     historiesCombined,
		VendorCheckList: vendorCheckList,
		Start:           lastPDFEndTime,
		End:             currentTime,
	})
	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

func (r *reportService) InsertPdf(input dto.PdfFile) (*string, rest_err.APIError) {
	currentTime := time.Now().Unix()
	if input.EndReportTime > currentTime {
		input.EndReportTime = currentTime
	}

	return r.dao.Pdf.InsertPdf(input)
}

func (r *reportService) FindPdf(branch string, typePdf string) ([]dto.PdfFile, rest_err.APIError) {
	return r.dao.Pdf.FindPdf(branch, typePdf)
}

func (r *reportService) GenerateReportVendorDaily(name string, branch string, target int64) (*string, rest_err.APIError) {
	currentTime := time.Now().Unix()
	if target > currentTime {
		target = currentTime
	}

	targetMinDaily := target - 60*60*24       // -1 hari
	targetMinMonthly := target - 60*60*24*60  // -2 bulan
	targetMinQuarter := target - 60*60*24*150 // -5 bulan

	// cek virtual cctv
	cctvVirtual, _ := r.dao.CheckCCTV.GetLastCheckCreateRange(targetMinDaily, target, branch)

	// cek fisik cctv bulanan
	cctvMonthly, _ := r.dao.CheckCCTVPhy.GetLastCheckCreateRange(targetMinMonthly, target, branch, false)

	// cek fisik cctv 3 bulanan
	cctvQuarter, _ := r.dao.CheckCCTVPhy.GetLastCheckCreateRange(targetMinQuarter, target, branch, true)

	// cek virtual altai
	altaiVirtual, _ := r.dao.CheckAltai.GetLastCheckCreateRange(targetMinDaily, target, branch)

	// cek fisik altai bulanan
	altaiMonthly, _ := r.dao.CheckAltaiPhy.GetLastCheckCreateRange(targetMinMonthly, target, branch, false)

	// cek fisik altai 3 bulanan
	altaiQuarter, _ := r.dao.CheckAltaiPhy.GetLastCheckCreateRange(targetMinQuarter, target, branch, true)

	// GET HISTORIES 0, 4 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: "0,4",
		}, dto.FilterTimeRangeLimit{
			FilterStart: targetMinDaily,
			FilterEnd:   target,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: "1,2,3",
		}, dto.FilterTimeRangeLimit{
			FilterStart: target - (3 * 30 * 24 * 60 * 60), // 3 bulan,
			FilterEnd:   target,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	historiesCombined := append(historyList04, historyList123...)

	errPDF := pdfgen.GeneratePDFVendorDaily(name, dto.ReportResponse{
		TargetTime:     target,
		CctvDaily:      cctvVirtual,
		CctvMonthly:    cctvMonthly,
		CctvQuarterly:  cctvQuarter,
		AltaiDaily:     altaiVirtual,
		AltaiMonthly:   altaiMonthly,
		AltaiQuarterly: altaiQuarter,
	}, pdfgen.PDFVendorReq{
		Name:        name,
		HistoryList: historiesCombined,
		Start:       targetMinDaily,
		End:         target,
	})

	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

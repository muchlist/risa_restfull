package service

import (
	"context"
	"fmt"
	"github.com/muchlist/risa_restfull/constants/enum"
	"github.com/muchlist/risa_restfull/dao/configcheckdao"
	"time"

	"github.com/muchlist/erru_utils_go/rest_err"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/pdftype"
	"github.com/muchlist/risa_restfull/dao/altaicheckdao"
	"github.com/muchlist/risa_restfull/dao/altaiphycheckdao"
	"github.com/muchlist/risa_restfull/dao/checkdao"
	"github.com/muchlist/risa_restfull/dao/historydao"
	"github.com/muchlist/risa_restfull/dao/reportdao"
	"github.com/muchlist/risa_restfull/dao/stockdao"
	"github.com/muchlist/risa_restfull/dao/vendorcheckdao"
	"github.com/muchlist/risa_restfull/dao/venphycheckdao"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/utils/pdfgen"
	"github.com/muchlist/risa_restfull/utils/pdfgen/stockpdf"
)

// ReportParams berisi semua dao yang diperlukan reports service, karena sangat banyak maka dibuat struct
type ReportParams struct {
	History       historydao.HistoryDaoAssumer
	CheckIT       checkdao.CheckDaoAssumer
	CheckCCTV     vendorcheckdao.CheckVendorDaoAssumer
	CheckCCTVPhy  venphycheckdao.CheckVenPhyDaoAssumer
	CheckAltai    altaicheckdao.CheckAltaiDaoAssumer
	CheckAltaiPhy altaiphycheckdao.CheckAltaiPhyDaoAssumer
	Stock         stockdao.StockDaoAssumer
	CheckConfig   configcheckdao.CheckConfigLoader
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
	GenerateReportVendorDaily(name string, branch string, start int64, end int64) (*string, rest_err.APIError)
	GenerateReportVendorDailyStartFromLast(name string, branch string) (*string, rest_err.APIError)
	GenerateReportPDFVendorMonthly(name string, branch string, start int64, end int64) (*string, rest_err.APIError)
	GenerateStockReportRestock(name, branch, category string, start, end int64) (*string, rest_err.APIError)
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

	// GET HISTORIES 0, 4, 6 sesuai start end inputan INFO, COMPLETE, COMPLETE BA
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       "",
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d", enum.HInfo, enum.HComplete, enum.HCompleteWithBA),
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
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d,%d", enum.HProgress, enum.HRequestPending, enum.HPending, enum.HRequestComplete),
		}, dto.FilterTimeRangeLimit{
			FilterStart: start - (3 * 30 * 24 * 60 * 60), // 3 bulan,
			FilterEnd:   end,
			Limit:       300,
		},
	)
	if err != nil {
		return nil, err
	}

	historiesCombined := append(historyList04, historyList123...)

	// GET CHECK LIST
	checkList, err := r.dao.CheckIT.FindCheckForReports(context.TODO(), branch, dto.FilterTimeRangeLimit{
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

	// GET HISTORIES 0, 4, 6 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       "",
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d", enum.HInfo, enum.HComplete, enum.HCompleteWithBA),
		}, dto.FilterTimeRangeLimit{
			FilterStart: lastPDFEndTime,
			FilterEnd:   currentTime,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3 , 5 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       "",
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d,%d", enum.HProgress, enum.HRequestPending, enum.HPending, enum.HRequestComplete),
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

	// GET CHECK LIST
	checkList, err := r.dao.CheckIT.FindCheckForReports(context.TODO(), branch, dto.FilterTimeRangeLimit{
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

	// GET HISTORIES 0, 4, 6 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d", enum.HInfo, enum.HComplete, enum.HCompleteWithBA),
		}, dto.FilterTimeRangeLimit{
			FilterStart: start,
			FilterEnd:   end,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3, 5 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d,%d", enum.HProgress, enum.HRequestPending, enum.HPending, enum.HRequestComplete),
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

	// GET HISTORIES 0, 4, 6 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d", enum.HInfo, enum.HComplete, enum.HCompleteWithBA),
		}, dto.FilterTimeRangeLimit{
			FilterStart: lastPDFEndTime,
			FilterEnd:   currentTime,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3, 5 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s", category.Cctv, category.Altai),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d,%d", enum.HProgress, enum.HRequestPending, enum.HPending, enum.HRequestComplete),
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

func (r *reportService) GenerateReportVendorDaily(name string, branch string, start int64, end int64) (*string, rest_err.APIError) {
	if start > end {
		return nil, rest_err.NewBadRequestError("tanggal awal tidak boleh lebih besar dari tanggal akhir")
	}

	currentTime := time.Now().Unix()
	if end > currentTime {
		end = currentTime
	}

	targetMinDaily := end - 60*60*24       // -1 hari
	targetMinMonthly := end - 60*60*24*60  // -2 bulan
	targetMinQuarter := end - 60*60*24*150 // -5 bulan

	if start > 0 {
		targetMinDaily = start
	}

	// cek virtual cctv
	cctvVirtual, _ := r.dao.CheckCCTV.GetLastCheckCreateRange(targetMinDaily, end, branch)

	// cek fisik cctv bulanan
	cctvMonthly, _ := r.dao.CheckCCTVPhy.GetLastCheckCreateRange(targetMinMonthly, end, branch, false)

	// cek fisik cctv 3 bulanan
	cctvQuarter, _ := r.dao.CheckCCTVPhy.GetLastCheckCreateRange(targetMinQuarter, end, branch, true)

	// cek virtual altai
	altaiVirtual, _ := r.dao.CheckAltai.GetLastCheckCreateRange(context.TODO(), targetMinDaily, end, branch)

	// cek fisik altai bulanan
	altaiMonthly, _ := r.dao.CheckAltaiPhy.GetLastCheckCreateRange(context.TODO(), targetMinMonthly, end, branch, false)

	// cek fisik altai 3 bulanan
	altaiQuarter, _ := r.dao.CheckAltaiPhy.GetLastCheckCreateRange(context.TODO(), targetMinQuarter, end, branch, true)

	// GET HISTORIES 0, 4, 6 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s,%s", category.Cctv, category.Altai, category.OtherV),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d", enum.HInfo, enum.HComplete, enum.HCompleteWithBA),
		}, dto.FilterTimeRangeLimit{
			FilterStart: targetMinDaily,
			FilterEnd:   end,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3, 5 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s,%s", category.Cctv, category.Altai, category.OtherV),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d,%d", enum.HProgress, enum.HRequestPending, enum.HPending, enum.HRequestComplete),
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

	errPDF := pdfgen.GeneratePDFVendorDaily(name, dto.ReportResponse{
		TargetTime:     end,
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
		End:         end,
	})

	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

func (r *reportService) GenerateReportVendorDailyStartFromLast(name string, branch string) (*string, rest_err.APIError) {
	currentTime := time.Now().Unix()
	lastPDF, err := r.dao.Pdf.FindLastPdf(branch, pdftype.Vendor)
	if err != nil {
		return nil, rest_err.NewBadRequestError("gagal mendapatkan data laporan sebelumnya")
	}
	lastPDFEndTime := lastPDF.EndReportTime

	if currentTime-lastPDFEndTime < 60*2 {
		return nil, rest_err.NewBadRequestError("Gagal. Jarak pembuatan laporan tidak boleh kurang dari 2 menit!")
	}

	targetMinDaily := currentTime - 60*60*24       // -1 hari
	targetMinMonthly := currentTime - 60*60*24*60  // -2 bulan
	targetMinQuarter := currentTime - 60*60*24*150 // -5 bulan

	// cek virtual cctv
	cctvVirtual, _ := r.dao.CheckCCTV.GetLastCheckCreateRange(targetMinDaily, currentTime, branch)

	// cek fisik cctv bulanan
	cctvMonthly, _ := r.dao.CheckCCTVPhy.GetLastCheckCreateRange(targetMinMonthly, currentTime, branch, false)

	// cek fisik cctv 3 bulanan
	cctvQuarter, _ := r.dao.CheckCCTVPhy.GetLastCheckCreateRange(targetMinQuarter, currentTime, branch, true)

	// cek virtual altai
	altaiVirtual, _ := r.dao.CheckAltai.GetLastCheckCreateRange(context.TODO(), targetMinDaily, currentTime, branch)

	// cek fisik altai bulanan
	altaiMonthly, _ := r.dao.CheckAltaiPhy.GetLastCheckCreateRange(context.TODO(), targetMinMonthly, currentTime, branch, false)

	// cek fisik altai 3 bulanan
	altaiQuarter, _ := r.dao.CheckAltaiPhy.GetLastCheckCreateRange(context.TODO(), targetMinQuarter, currentTime, branch, true)

	// GET HISTORIES 0, 4, 6 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s,%s", category.Cctv, category.Altai, category.OtherV),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d", enum.HInfo, enum.HComplete, enum.HCompleteWithBA),
		}, dto.FilterTimeRangeLimit{
			FilterStart: lastPDFEndTime,
			FilterEnd:   currentTime,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	// GET HISTORIES 1, 2, 3, 5 sesuai end inputan dan start = end - 3 bulan
	historyList123, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s, %s", category.Cctv, category.Altai, category.OtherV),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d,%d", enum.HProgress, enum.HRequestPending, enum.HPending, enum.HRequestComplete),
		}, dto.FilterTimeRangeLimit{
			FilterStart: currentTime - (3 * 30 * 24 * 60 * 60), // 3 bulan,
			FilterEnd:   currentTime,
			Limit:       300,
		},
	)

	if err != nil {
		return nil, err
	}

	historiesCombined := append(historyList04, historyList123...)

	errPDF := pdfgen.GeneratePDFVendorDaily(name, dto.ReportResponse{
		TargetTime:     currentTime,
		CctvDaily:      cctvVirtual,
		CctvMonthly:    cctvMonthly,
		CctvQuarterly:  cctvQuarter,
		AltaiDaily:     altaiVirtual,
		AltaiMonthly:   altaiMonthly,
		AltaiQuarterly: altaiQuarter,
	}, pdfgen.PDFVendorReq{
		Name:        name,
		HistoryList: historiesCombined,
		Start:       lastPDFEndTime,
		End:         currentTime,
	})

	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

// GenerateReportPDFVendorMonthly membuat Pdf untuk vendor multinet
func (r *reportService) GenerateReportPDFVendorMonthly(name string, branch string, start int64, end int64) (*string, rest_err.APIError) {
	if start > end && start < 0 {
		return nil, rest_err.NewBadRequestError("tanggal awal tidak boleh lebih besar dari tanggal akhir")
	}

	currentTime := time.Now().Unix()
	if end > currentTime {
		end = currentTime
	}

	targetMinMonthly := end - 60*60*24*60  // -2 bulan
	targetMinQuarter := end - 60*60*24*150 // -5 bulan

	// GET HISTORIES 4 sesuai start end inputan
	historyList04, err := r.dao.History.UnwindHistory(
		dto.FilterBranchCatInCompleteIn{
			FilterBranch:         branch,
			FilterCategory:       fmt.Sprintf("%s,%s,%s", category.Cctv, category.Altai, category.OtherV),
			FilterCompleteStatus: fmt.Sprintf("%d,%d", enum.HComplete, enum.HCompleteWithBA),
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
			FilterCategory:       fmt.Sprintf("%s,%s,%s", category.Cctv, category.Altai, category.OtherV),
			FilterCompleteStatus: fmt.Sprintf("%d,%d,%d,%d", enum.HProgress, enum.HRequestPending, enum.HPending, enum.HRequestComplete),
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

	// cek fisik cctv bulanan
	cctvMonthly, _ := r.dao.CheckCCTVPhy.GetLastCheckCreateRange(targetMinMonthly, currentTime, branch, false)

	// cek fisik cctv 3 bulanan
	cctvQuarter, _ := r.dao.CheckCCTVPhy.GetLastCheckCreateRange(targetMinQuarter, currentTime, branch, true)

	// cek fisik altai bulanan
	altaiMonthly, _ := r.dao.CheckAltaiPhy.GetLastCheckCreateRange(context.TODO(), targetMinMonthly, currentTime, branch, false)

	// cek fisik altai 3 bulanan
	altaiQuarter, _ := r.dao.CheckAltaiPhy.GetLastCheckCreateRange(context.TODO(), targetMinQuarter, currentTime, branch, true)

	// checklist backup config
	lastCheckConfig, _ := r.dao.CheckConfig.GetLastCheckCreateRange(context.TODO(), targetMinMonthly, currentTime, branch)

	errPDF := pdfgen.GeneratePDFVendorMonthly(pdfgen.PDFReqMonth{
		Name:        name,
		HistoryList: historiesCombined,
		Start:       start,
		End:         end,
	}, dto.ReportResponse{
		TargetTime:     end,
		CctvMonthly:    cctvMonthly,
		CctvQuarterly:  cctvQuarter,
		AltaiMonthly:   altaiMonthly,
		AltaiQuarterly: altaiQuarter,
	},
		*lastCheckConfig,
	)
	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

// GenerateStockReportRestock membuat Pdf untuk stock yang perlu diisi ulang
func (r *reportService) GenerateStockReportRestock(name, branch, category string, start, end int64) (*string, rest_err.APIError) {
	if start > end && start < 0 {
		return nil, rest_err.NewBadRequestError("tanggal awal tidak boleh lebih besar dari tanggal akhir")
	}

	currentTime := time.Now().Unix()
	if end > currentTime {
		end = currentTime
	}

	// GET Stock need Restock
	stockList, err := r.dao.Stock.FindStockNeedRestock(
		dto.FilterBranchCatDisable{
			FilterBranch:   branch,
			FilterCategory: category,
			FilterDisable:  false,
		},
	)
	if err != nil {
		return nil, err
	}

	errPDF := stockpdf.GenerateStockPDF(stockpdf.PDFReq{
		Name:      name,
		StockList: stockList,
		Start:     start,
		End:       end,
	})
	if errPDF != nil {
		return nil, rest_err.NewInternalServerError("gagal membuat Pdf", errPDF)
	}

	return &name, nil
}

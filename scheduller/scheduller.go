package scheduller

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/constants/pdftype"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/muchlist/risa_restfull/utils/timegen"
)

func RunScheduler(
	genUnitService service.GenUnitServiceAssumer,
	reportService service.ReportServiceAssumer,
) {
	witaTimeZone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		logger.Error("gagal menggunakan timezone wita", err)
	}
	s := gocron.NewScheduler(witaTimeZone)

	// run check cctv
	_, _ = s.Every(2).Hours().StartAt(time.Now().Add(time.Second * 30)).Do(func() {
		runCctvCheckBanjarmasin(genUnitService)
	})

	// run report vendor berjalan setiap jam 1 malam tanggal 1 bulan berjalan
	_, _ = s.Every(1).Month(1).At("01:00").Do(func() {
		runReportGeneratorVendormonthlyBanjarmasin(reportService)
	})

	s.StartAsync()
}

func runCctvCheckBanjarmasin(genUnitService service.GenUnitServiceAssumer) {
	_ = genUnitService.CheckHardwareDownAndSendNotif(context.Background(), "BANJARMASIN", category.Cctv)
}

func runReportGeneratorVendormonthlyBanjarmasin(reportService service.ReportServiceAssumer) {

	// berjalan setiap tanggal 1 bulan sekarang jam 00.01
	// cari tanggal terakhir bulan sebelumnya jam 24.00 --- >  time end
	// cari tanggal awal bulan sebelumnya jam 00.00   --- >  time start

	loc, _ := time.LoadLocation("Asia/Makassar")
	timeNow := time.Now().In(loc)
	month := timeNow.Month()
	beforeMonth := int(month) - 1
	if beforeMonth == 0 {
		beforeMonth = 12
	}

	// tanggal awal
	timeStart := time.Date(timeNow.Year(), time.Month(beforeMonth), 1, 00, 00, 00, 0, loc)
	timeEnd := timeStart.AddDate(0, 1, 0).Add(time.Second * -1)

	timeStartUnix := timeStart.Unix()
	timeEndUnix := timeEnd.Unix()

	pdfName, err := timegen.GetTimeAsName(timeEndUnix)
	if err != nil {
		logger.Error("gagal membuat nama pdf", err)
	}
	pdfName = fmt.Sprintf("vendor-monthly%s", pdfName)

	_, apiErr := reportService.GenerateReportPDFVendorMonthly(context.Background(), pdfName, "BANJARMASIN", timeStartUnix, timeEndUnix, false)
	if apiErr != nil {
		logger.Error(apiErr.Message(), apiErr)
	}

	// simpan pdf ke database
	_, apiErr = reportService.InsertPdf(context.Background(), dto.PdfFile{
		CreatedAt:     time.Now().Unix(),
		CreatedBy:     "SYSTEM",
		Branch:        "BANJARMASIN",
		Name:          pdfName,
		Type:          pdftype.VendorMonthly,
		FileName:      fmt.Sprintf("pdf-v-month/%s.pdf", pdfName),
		EndReportTime: timeEndUnix,
	})

	if apiErr != nil {
		logger.Error(apiErr.Message(), apiErr)
	}
}

/*
func runReportGeneratorBanjarmasin(reportService service.ReportServiceAssumer) {
	timeNow := time.Now().Unix()
	timePast := timeNow - 28801 // minus 8 jam

	pdfName, err2 := timegen.GetTimeAsName(timeNow)
	if err2 != nil {
		logger.Error("error membuat nama pdf otomatis", err2)
	}
	pdfName = fmt.Sprintf("auto-%s", pdfName)

	_, err := reportService.GenerateReportPDF(pdfName, "BANJARMASIN", timePast, timeNow)
	if err != nil {
		logger.Error("gagal membuat pdf otomatis", err)
	}

	// simpan pdf ke database
	_, apiErr := reportService.InsertPdf(dto.PdfFile{
		CreatedAt: time.Now().Unix(),
		CreatedBy: "SYSTEM",
		Branch:    "BANJARMASIN",
		Name:      pdfName,
		Type:      "LAPORAN",
		FileName:  fmt.Sprintf("pdf/%s.pdf", pdfName),
	})

	if apiErr != nil {
		logger.Error("gagal membuat pdf otomatis", err)
	}
}*/

/*
func runSpeedTest(speedService service.SpeedTestServiceAssumer) {
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		logger.Error("speedTest gagal dijalankan", err)
		return
	}

	serverList, err := speedtest.FetchServerList(user)
	if err != nil {
		logger.Error("speedTest gagal dijalankan", err)
		return
	}

	targets, err := serverList.FindServer([]int{})
	if err != nil {
		logger.Error("speedTest gagal dijalankan", err)
		return
	}

	for _, s := range targets {
		_ = s.PingTest()
		_ = s.DownloadTest(false)
		_ = s.UploadTest(false)

		data := dto.SpeedTest{
			LatencyMs: s.Latency.Milliseconds(),
			Download:  s.DLSpeed,
			Upload:    s.ULSpeed,
		}
		_, err := speedService.InsertSpeed(data)
		if err != nil {
			logger.Error("speedTest gagal dijalankan (insertSpeed ke database)", err)
			return
		}
	}
}
*/

package scheduller

import (
	"github.com/go-co-op/gocron"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/service"
	"time"
)

func RunScheduler(
	genUnitService service.GenUnitServiceAssumer,
	//reportService service.ReportServiceAssumer,
) {
	witaTimeZone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		logger.Error("gagal menggunakan timezone wita", err)
	}
	s := gocron.NewScheduler(witaTimeZone)

	/*run speed test
	_, _ = s.Every(1).Days().At("06:00").Do(func() {
		runSpeedTest(speedService)
	})*/

	// run check cctv
	_, _ = s.Every(2).Hours().StartAt(time.Now().Add(time.Second * 30)).Do(func() {
		runCctvCheckBanjarmasin(genUnitService)
	})

	// run report vendor
	/*_, _ = s.Every(1).Days().At("17:00").Do(func() {
		runReportGeneratorVendorBanjarmasin(reportService)
	})*/

	// run report generator
	/*	_, _ = s.Every(1).Days().At("08:00").Do(func() {
			runReportGeneratorBanjarmasin(reportService)
		})
		_, _ = s.Every(1).Days().At("16:00").Do(func() {
			runReportGeneratorBanjarmasin(reportService)
		})

		_, _ = s.Every(1).Days().At("00:00").Do(func() {
			runReportGeneratorBanjarmasin(reportService)
		})*/

	s.StartAsync()
}

func runCctvCheckBanjarmasin(genUnitService service.GenUnitServiceAssumer) {
	_ = genUnitService.CheckHardwareDownAndSendNotif("BANJARMASIN", category.Cctv)
}

/*func runReportGeneratorVendorBanjarmasin(reportService service.ReportServiceAssumer) {
	timeNow := time.Now().Unix()
	timePast := timeNow - (60 * 60 * 24) // minus 24 jam

	pdfName, err2 := timegen.GetTimeAsName(timeNow)
	if err2 != nil {
		logger.Error("error membuat nama pdf otomatis", err2)
	}
	pdfName = fmt.Sprintf("auto-vendor-%s", pdfName)

	_, err := reportService.GenerateReportPDFVendor(pdfName, "BANJARMASIN", timePast, timeNow)
	if err != nil {
		logger.Error("gagal membuat pdf otomatis", err)
	}

	// simpan pdf ke database
	_, apiErr := reportService.InsertPdf(dto.PdfFile{
		CreatedAt: time.Now().Unix(),
		CreatedBy: "SYSTEM",
		Branch:    "BANJARMASIN",
		Name:      pdfName,
		Type:      "VENDOR",
		FileName:  fmt.Sprintf("pdf-vendor/%s.pdf", pdfName),
	})

	if apiErr != nil {
		logger.Error("gagal membuat pdf otomatis", err)
	}
}*/

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

package scheduller

import (
	"github.com/go-co-op/gocron"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/risa_restfull/constants/category"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/showwin/speedtest-go/speedtest"
	"time"
)

func RunScheduler(
	speedService service.SpeedTestServiceAssumer,
	genUnitService service.GenUnitServiceAssumer,
	reportService service.ReportServiceAssumer,
) {
	witaTimeZone, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		logger.Error("gagal menggunakan timezone wita", err)
	}
	s := gocron.NewScheduler(witaTimeZone)

	// run speed test
	_, _ = s.Every(1).Days().At("23:00").Do(func() {
		runSpeedTest(speedService)
	})

	// run check cctv
	_, _ = s.Every(2).Hours().Do(func() {
		runCctvCheckBanjarmasin(genUnitService)
	})

	// run report generator
	_, _ = s.Every(1).Days().At("08:00").Do(func() {
		runReportGeneratorBanjarmasin(reportService)
	})
	_, _ = s.Every(1).Days().At("16:00").Do(func() {
		runReportGeneratorBanjarmasin(reportService)
	})
	_, _ = s.Every(1).Days().At("00:00").Do(func() {
		runReportGeneratorBanjarmasin(reportService)
	})

	s.StartAsync()
}

func runCctvCheckBanjarmasin(genUnitService service.GenUnitServiceAssumer) {
	_ = genUnitService.CheckHardwareDownAndSendNotif("BANJARMASIN", category.Cctv)
}

func runReportGeneratorBanjarmasin(reportService service.ReportServiceAssumer) {
	timeNow := time.Now().Unix()
	timePast := timeNow - 28801 // minus 8 jam
	_, err := reportService.GenerateReportPDF("BANJARMASIN", timePast, timeNow)
	if err != nil {
		logger.Error("gagal membuat pdf otomatis", err)
	}
}

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

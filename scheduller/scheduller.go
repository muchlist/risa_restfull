package scheduller

import (
	"github.com/go-co-op/gocron"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/risa_restfull/dto"
	"github.com/muchlist/risa_restfull/service"
	"github.com/showwin/speedtest-go/speedtest"
	"time"
)

func RunScheduler(speedService service.SpeedTestServiceAssumer) {
	s := gocron.NewScheduler(time.UTC)
	_, _ = s.Every(4).Hour().Do(func() {
		runSpeedTest(speedService)
	})
	s.StartAsync()
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

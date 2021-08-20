package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/erru_utils_go/logger"
	"github.com/muchlist/risa_restfull/clients/fcm"
	"github.com/muchlist/risa_restfull/db"
	"github.com/muchlist/risa_restfull/scheduller"
	"github.com/muchlist/risa_restfull/utils/mjwt"
)

func RunApp() {
	// inisiasi database mongodb
	client, ctx, cancel := db.Init()
	defer client.Disconnect(ctx) //nolint:errcheck
	defer cancel()

	// inisasi firebase app
	_ = fcm.Init()

	// inisiasi jwt
	mjwt.Init()

	app := fiber.New()
	// memenuhi dependency, mapping url
	mapUrls(app)

	// menjalankan job scheduller cctv
	scheduller.RunScheduler(genUnitService)

	if err := app.Listen(":3500"); err != nil {
		logger.Error("error fiber listen", err)
		return
	}
}

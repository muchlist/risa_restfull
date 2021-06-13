package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/risa_restfull/clients/fcm"
	"github.com/muchlist/risa_restfull/db"
	"log"
)

func RunApp() {
	// inisiasi database mongodb
	client, ctx, cancel := db.Init()
	defer client.Disconnect(ctx) //nolint:errcheck
	defer cancel()

	// inisasi firebase app
	_ = fcm.Init()

	app := fiber.New()
	mapUrls(app)
	if err := app.Listen(":3500"); err != nil {
		log.Print(err)
		return
	}
}

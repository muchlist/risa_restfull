package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/risa_restfull/db"
	"log"
)

func RunApp() {

	// inisiasi database
	client, ctx, cancel := db.Init()
	defer client.Disconnect(ctx) //nolint:errcheck
	defer cancel()

	app := fiber.New()
	mapUrls(app)
	if err := app.Listen(":3500"); err != nil {
		log.Print(err)
		return
	}
}

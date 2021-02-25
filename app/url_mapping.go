package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/middleware"
)

func mapUrls(app *fiber.App) {
	app.Use(logger.New())
	app.Use(middleware.LimitRequest())

	app.Static("/image/avatar", "./static/image/avatar")
	app.Static("/image/history", "./static/image/history")

	api := app.Group("/api/v1")
	api.Get("/ping", pingHandler.Ping)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)

	api.Get("/users", middleware.NormalAuth(), userHandler.Find)
	api.Get("/profile", middleware.NormalAuth(), userHandler.GetProfile)
	api.Post("/avatar", middleware.NormalAuth(), userHandler.UploadImage)

	api.Post("/change-password", middleware.FreshAuth(), userHandler.ChangePassword)

	apiAuthAdmin := app.Group("/api/v1/admin")
	apiAuthAdmin.Use(middleware.NormalAuth(roles.RoleAdmin))
	apiAuthAdmin.Post("/users", userHandler.Register)
	apiAuthAdmin.Put("/users/:user_id", userHandler.Edit)
	apiAuthAdmin.Delete("/users/:user_id", userHandler.Delete)
	apiAuthAdmin.Get("/users/:user_id/reset-password", userHandler.ResetPassword)

	//Unit GENERAL
	api.Get("/general", middleware.NormalAuth(), genUnitHandler.Find)

	//History
	api.Get("/histories", middleware.NormalAuth(), historyHandler.Find)
	api.Post("/histories", middleware.NormalAuth(), historyHandler.Insert)
}

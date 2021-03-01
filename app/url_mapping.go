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
	app.Static("/image/cctv", "./static/image/cctv")

	api := app.Group("/api/v1")

	//PING
	api.Get("/ping", pingHandler.Ping)

	//USER
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/users", middleware.NormalAuth(), userHandler.Find)
	api.Get("/profile", middleware.NormalAuth(), userHandler.GetProfile)
	api.Post("/avatar", middleware.NormalAuth(), userHandler.UploadImage)
	api.Post("/change-password", middleware.FreshAuth(), userHandler.ChangePassword)

	//USER ADMIN
	apiAuthAdmin := app.Group("/api/v1/admin")
	apiAuthAdmin.Use(middleware.NormalAuth(roles.RoleAdmin))
	apiAuthAdmin.Post("/users", userHandler.Register)
	apiAuthAdmin.Put("/users/:user_id", userHandler.Edit)
	apiAuthAdmin.Delete("/users/:user_id", userHandler.Delete)
	apiAuthAdmin.Get("/users/:user_id/reset-password", userHandler.ResetPassword)

	//Unit GENERAL
	api.Get("/general", middleware.NormalAuth(), genUnitHandler.Find)
	api.Get("/general-ip", genUnitHandler.GetIPList)
	api.Post("/general-ip-state", genUnitHandler.UpdatePingState)

	//History
	api.Get("/histories", middleware.NormalAuth(), historyHandler.Find)
	api.Get("/histories/:id", middleware.NormalAuth(), historyHandler.GetHistory)
	api.Delete("/histories/:id", middleware.NormalAuth(), historyHandler.Delete)
	api.Get("/histories-parent/:id", middleware.NormalAuth(), historyHandler.FindFromParent)
	api.Get("/histories-user/:id", middleware.NormalAuth(), historyHandler.FindFromUser)
	api.Post("/histories", middleware.NormalAuth(), historyHandler.Insert)
	api.Put("/histories/:id", middleware.NormalAuth(), historyHandler.Edit)
	api.Post("/history-image/:id", middleware.NormalAuth(), historyHandler.UploadImage) // IMPROVEMENT post image when build history

	//CCTV
	api.Post("/cctv", middleware.NormalAuth(), cctvHandler.Insert)
	api.Get("/cctv/:id", middleware.NormalAuth(), cctvHandler.GetCctv)
	api.Put("/cctv/:id", middleware.NormalAuth(), cctvHandler.Edit)
	api.Delete("/cctv/:id", middleware.NormalAuth(), cctvHandler.Delete)
	api.Get("/cctv", middleware.NormalAuth(), cctvHandler.Find) // IMPROVEMENT join table with gen_unit
	api.Get("/cctv-avail/:id/:status", middleware.NormalAuth(), cctvHandler.DisableCctv)
	api.Post("/cctv-image/:id", middleware.NormalAuth(), cctvHandler.UploadImage)

	//STOCK
	api.Post("/stock", middleware.NormalAuth(), stockHandler.Insert)
	api.Get("/stock/:id", middleware.NormalAuth(), stockHandler.GetStock)
	api.Put("/stock/:id", middleware.NormalAuth(), stockHandler.Edit)
	api.Delete("/stock/:id", middleware.NormalAuth(), stockHandler.Delete)
	api.Post("/stock-change/:id", middleware.NormalAuth(), stockHandler.ChangeQty)
	api.Get("/stock", middleware.NormalAuth(), stockHandler.Find)
	api.Get("/stock-avail/:id/:status", middleware.NormalAuth(), stockHandler.DisableStock)
	api.Post("/stock-image/:id", middleware.NormalAuth(), stockHandler.UploadImage)

	//CHECK ITEM
	api.Post("/check-item", middleware.NormalAuth(), checkItemHandler.Insert)
	api.Get("/check-item/:id", middleware.NormalAuth(), checkItemHandler.GetCheckItem)
	api.Put("/check-item/:id", middleware.NormalAuth(), checkItemHandler.Edit)
	api.Delete("/check-item/:id", middleware.NormalAuth(), checkItemHandler.Delete)
	api.Get("/check-item", middleware.NormalAuth(), checkItemHandler.Find)
	api.Get("/check-item-avail/:id/:status", middleware.NormalAuth(), checkItemHandler.DisableCheckItem)

	//CHECK
	api.Post("/check", middleware.NormalAuth(), checkHandler.Insert)
	api.Get("/check/:id", middleware.NormalAuth(), checkHandler.GetCheck)
	api.Put("/check/:id", middleware.NormalAuth(), checkHandler.Edit)
	api.Delete("/check/:id", middleware.NormalAuth(), checkHandler.Delete)
	api.Get("/check", middleware.NormalAuth(), checkHandler.Find)
	api.Post("/check-update/:id/:child_id", middleware.NormalAuth(), checkHandler.UpdateCheckItem)
	api.Post("/check-image/:id/:child_id", middleware.NormalAuth(), checkHandler.UploadImage)

}

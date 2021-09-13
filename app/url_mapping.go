package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/muchlist/risa_restfull/constants/roles"
	"github.com/muchlist/risa_restfull/handler"
	"github.com/muchlist/risa_restfull/middleware"
)

//nolint:funlen
func mapUrls(app *fiber.App) {

	// Controller or Handler
	pingHandler := handler.NewPingHandler()
	optionHandler := handler.NewOptionHandler()
	userHandler := handler.NewUserHandler(userService)
	genUnitHandler := handler.NewGenUnitHandler(genUnitService)
	historyHandler := handler.NewHistoryHandler(historyService)
	cctvHandler := handler.NewCctvHandler(cctvService)
	stockHandler := handler.NewStockHandler(stockService)
	checkItemHandler := handler.NewCheckItemHandler(checkItemService)
	checkHandler := handler.NewCheckHandler(checkService)
	improveHandler := handler.NewImproveHandler(improveService)
	computerHandler := handler.NewComputerHandler(computerService)
	otherHandler := handler.NewOtherHandler(otherService)
	vendorCheckHandler := handler.NewVendorCheckHandler(vendorCheckService)
	altaiCheckHandler := handler.NewAltaiCheckHandler(altaiCheckService)
	venPhyCheckHandler := handler.NewVenPhyCheckHandler(venPhyCheckService)
	altaiPhyCheckHandler := handler.NewAltaiPhyCheckHandler(altaiPhyCheckService)
	speedHandler := handler.NewSpeedHandler(speedService)
	reportHandler := handler.NewReportHandler(reportService)

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type, Accept, Authorization",
	}))
	app.Use(middleware.LimitRequest())

	app.Static("/image/avatar", "./static/image/avatar")
	app.Static("/image/history", "./static/image/history")
	app.Static("/image/cctv", "./static/image/cctv")
	app.Static("/image/computer", "./static/image/computer")
	app.Static("/image/other", "./static/image/other")
	app.Static("/image/stock", "./static/image/stock")
	app.Static("/image/check", "./static/image/check")
	app.Static("/image/config", "./static/image/config")
	app.Static("/pdf", "./static/pdf")
	app.Static("/pdf-vendor", "./static/pdf-vendor")
	app.Static("/pdf-v-month", "./static/pdf-v-month")
	app.Static("/pdf-stock", "./static/pdf-stock")

	api := app.Group("/api/v1")

	// PING
	api.Get("/ping", pingHandler.Ping)

	// USER
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/users", middleware.NormalAuth(), userHandler.Find)
	api.Get("/profile", middleware.NormalAuth(), userHandler.GetProfile)
	api.Post("/avatar", middleware.NormalAuth(), userHandler.UploadImage)
	api.Post("/change-password", middleware.FreshAuth(), userHandler.ChangePassword)
	api.Post("/update-fcm", middleware.NormalAuth(), userHandler.UpdateFcmToken)

	// USER ADMIN
	apiAuthAdmin := app.Group("/api/v1/admin")
	apiAuthAdmin.Use(middleware.NormalAuth(roles.RoleAdmin))
	apiAuthAdmin.Post("/users", userHandler.Register)
	apiAuthAdmin.Put("/users/:user_id", userHandler.Edit)
	apiAuthAdmin.Delete("/users/:user_id", userHandler.Delete)
	apiAuthAdmin.Get("/users/:user_id/reset-password", userHandler.ResetPassword)

	// Unit GENERAL
	api.Get("/general", middleware.NormalAuth(), genUnitHandler.Find)
	api.Get("/general-ip", genUnitHandler.GetIPList)
	api.Post("/general-ip-state", genUnitHandler.UpdatePingState)

	// History
	api.Get("/histories", middleware.NormalAuth(), historyHandler.Find)
	api.Get("/histories-unwind", middleware.NormalAuth(), historyHandler.FindUnwind)
	api.Get("/histories/:id", middleware.NormalAuth(), historyHandler.GetHistory)
	api.Delete("/histories/:id", middleware.NormalAuth(), historyHandler.Delete)
	api.Get("/histories-parent/:id", middleware.NormalAuth(), historyHandler.FindFromParent)
	api.Get("/histories-user/:id", middleware.NormalAuth(), historyHandler.FindFromUser)
	api.Post("/histories", middleware.NormalAuth(), historyHandler.Insert)
	api.Put("/histories/:id", middleware.NormalAuth(), historyHandler.Edit)
	api.Post("/history-image/:id", middleware.NormalAuth(), historyHandler.UploadImage)
	api.Post("/upload-image/", middleware.NormalAuth(), historyHandler.UploadImageWithoutParent)

	// CCTV
	api.Post("/cctv", middleware.NormalAuth(), cctvHandler.Insert)
	api.Get("/cctv/:id", middleware.NormalAuth(), cctvHandler.GetCctv)
	api.Put("/cctv/:id", middleware.NormalAuth(), cctvHandler.Edit)
	api.Delete("/cctv/:id", middleware.NormalAuth(), cctvHandler.Delete)
	api.Get("/cctv", middleware.NormalAuth(), cctvHandler.Find)
	api.Get("/cctv-avail/:id/:status", middleware.NormalAuth(), cctvHandler.DisableCctv)
	api.Post("/cctv-image/:id", middleware.NormalAuth(), cctvHandler.UploadImage)

	// COMPUTER
	api.Post("/computer", middleware.NormalAuth(), computerHandler.Insert)
	api.Get("/computer/:id", middleware.NormalAuth(), computerHandler.GetComputer)
	api.Put("/computer/:id", middleware.NormalAuth(), computerHandler.Edit)
	api.Delete("/computer/:id", middleware.NormalAuth(), computerHandler.Delete)
	api.Get("/computer", middleware.NormalAuth(), computerHandler.Find)
	api.Get("/computer-avail/:id/:status", middleware.NormalAuth(), computerHandler.DisableComputer)
	api.Post("/computer-image/:id", middleware.NormalAuth(), computerHandler.UploadImage)

	// OTHER
	api.Post("/other", middleware.NormalAuth(), otherHandler.Insert)
	api.Get("/other/:id", middleware.NormalAuth(), otherHandler.GetOther)
	api.Put("/other/:id", middleware.NormalAuth(), otherHandler.Edit)
	api.Delete("/other/:cat/:id", middleware.NormalAuth(), otherHandler.Delete)
	api.Get("/others/:cat", middleware.NormalAuth(), otherHandler.Find)
	api.Get("/other-avail/:cat/:id/:status", middleware.NormalAuth(), otherHandler.DisableOther)
	api.Post("/other-image/:id", middleware.NormalAuth(), otherHandler.UploadImage)

	// STOCK
	api.Post("/stock", middleware.NormalAuth(), stockHandler.Insert)
	api.Get("/stock/:id", middleware.NormalAuth(), stockHandler.GetStock)
	api.Put("/stock/:id", middleware.NormalAuth(), stockHandler.Edit)
	api.Delete("/stock/:id", middleware.NormalAuth(), stockHandler.Delete)
	api.Post("/stock-change/:id", middleware.NormalAuth(), stockHandler.ChangeQty)
	api.Get("/stock", middleware.NormalAuth(), stockHandler.Find)
	api.Get("/restock", middleware.NormalAuth(), stockHandler.FindNeedRestock1)
	api.Get("/restock-2", middleware.NormalAuth(), stockHandler.FindNeedRestock2)
	api.Get("/stock-avail/:id/:status", middleware.NormalAuth(), stockHandler.DisableStock)
	api.Post("/stock-image/:id", middleware.NormalAuth(), stockHandler.UploadImage)

	// CHECK ITEM
	api.Post("/check-item", middleware.NormalAuth(), checkItemHandler.Insert)
	api.Get("/check-item/:id", middleware.NormalAuth(), checkItemHandler.GetCheckItem)
	api.Put("/check-item/:id", middleware.NormalAuth(), checkItemHandler.Edit)
	api.Delete("/check-item/:id", middleware.NormalAuth(), checkItemHandler.Delete)
	api.Get("/check-item", middleware.NormalAuth(), checkItemHandler.Find)
	api.Get("/check-item-avail/:id/:status", middleware.NormalAuth(), checkItemHandler.DisableCheckItem)

	// CHECK
	api.Post("/check", middleware.NormalAuth(), checkHandler.Insert)
	api.Get("/check/:id", middleware.NormalAuth(), checkHandler.GetCheck)
	api.Put("/check/:id", middleware.NormalAuth(), checkHandler.Edit)
	api.Delete("/check/:id", middleware.NormalAuth(), checkHandler.Delete)
	api.Get("/check", middleware.NormalAuth(), checkHandler.Find)
	api.Post("/check-update", middleware.NormalAuth(), checkHandler.UpdateCheckItem)
	api.Post("/check-image/:id/:child_id", middleware.NormalAuth(), checkHandler.UploadImage)

	// CCTV CHECK VIRTUAL
	api.Post("/vendor-check", middleware.NormalAuth(), vendorCheckHandler.Insert)
	api.Delete("/vendor-check/:id", middleware.NormalAuth(), vendorCheckHandler.Delete)
	api.Get("/vendor-check/:id", middleware.NormalAuth(), vendorCheckHandler.Get)
	api.Get("/vendor-check", middleware.NormalAuth(), vendorCheckHandler.Find)
	api.Post("/vendor-check-update", middleware.NormalAuth(), vendorCheckHandler.UpdateCheckItem)
	api.Post("/bulk-vendor-update", middleware.NormalAuth(), vendorCheckHandler.BulkUpdateCheckItem)
	api.Get("/vendor-check-finish/:id", middleware.NormalAuth(), vendorCheckHandler.Finish)

	// ALTAI CHECK VIRTUAL
	api.Post("/altai-check", middleware.NormalAuth(), altaiCheckHandler.Insert)
	api.Delete("/altai-check/:id", middleware.NormalAuth(), altaiCheckHandler.Delete)
	api.Get("/altai-check/:id", middleware.NormalAuth(), altaiCheckHandler.Get)
	api.Get("/altai-check", middleware.NormalAuth(), altaiCheckHandler.Find)
	api.Post("/altai-check-update", middleware.NormalAuth(), altaiCheckHandler.UpdateCheckItem)
	api.Post("/bulk-altai-update", middleware.NormalAuth(), altaiCheckHandler.BulkUpdateCheckItem)
	api.Get("/altai-check-finish/:id", middleware.NormalAuth(), altaiCheckHandler.Finish)

	// CCTV CHECK FISIK
	api.Post("/phy-check", middleware.NormalAuth(), venPhyCheckHandler.Insert)
	api.Post("/phy-check-quarter", middleware.NormalAuth(), venPhyCheckHandler.InsertQuarter)
	api.Delete("/phy-check/:id", middleware.NormalAuth(), venPhyCheckHandler.Delete)
	api.Get("/phy-check/:id", middleware.NormalAuth(), venPhyCheckHandler.Get)
	api.Get("/phy-check", middleware.NormalAuth(), venPhyCheckHandler.Find)
	api.Post("/phy-check-update", middleware.NormalAuth(roles.RoleVendor), venPhyCheckHandler.UpdateCheckItem)
	api.Post("/bulk-phy-update", middleware.NormalAuth(roles.RoleVendor), venPhyCheckHandler.BulkUpdateCheckItem)
	api.Get("/phy-check-finish/:id", middleware.NormalAuth(), venPhyCheckHandler.Finish)

	// ALTAI CHECK FISIK
	api.Post("/altai-phy-check", middleware.NormalAuth(), altaiPhyCheckHandler.Insert)
	api.Post("/altai-phy-check-quarter", middleware.NormalAuth(), altaiPhyCheckHandler.InsertQuarter)
	api.Delete("/altai-phy-check/:id", middleware.NormalAuth(), altaiPhyCheckHandler.Delete)
	api.Get("/altai-phy-check/:id", middleware.NormalAuth(), altaiPhyCheckHandler.Get)
	api.Get("/altai-phy-check", middleware.NormalAuth(), altaiPhyCheckHandler.Find)
	api.Post("/altai-phy-check-update", middleware.NormalAuth(roles.RoleVendor), altaiPhyCheckHandler.UpdateCheckItem)
	api.Post("/altai-bulk-phy-update", middleware.NormalAuth(roles.RoleVendor), altaiPhyCheckHandler.BulkUpdateCheckItem)
	api.Get("/altai-phy-check-finish/:id", middleware.NormalAuth(), altaiPhyCheckHandler.Finish)

	// IMPROVE
	api.Post("/improve", middleware.NormalAuth(), improveHandler.Insert)
	api.Get("/improve/:id", middleware.NormalAuth(), improveHandler.GetImprove)
	api.Put("/improve/:id", middleware.NormalAuth(roles.RoleApprove), improveHandler.Edit)
	api.Delete("/improve/:id", middleware.NormalAuth(), improveHandler.Delete)
	api.Post("/improve-change/:id", middleware.NormalAuth(), improveHandler.ChangeImprove)
	api.Get("/improve", middleware.NormalAuth(), improveHandler.Find)
	api.Get("/improve-status/:id/:status", middleware.NormalAuth(roles.RoleApprove), improveHandler.ActivateImprove)

	// speed test inet
	api.Get("/speed-test", middleware.NormalAuth(), speedHandler.Retrieve)

	// REPORT
	api.Get("/generate-pdf", middleware.NormalAuth(), reportHandler.GeneratePDF)
	api.Get("/generate-pdf-auto", middleware.NormalAuth(), reportHandler.GeneratePDFStartFromLast)
	api.Get("/generate-pdf-vendor", middleware.NormalAuth(), reportHandler.GeneratePDFVendor)
	api.Get("/generate-pdf-vendor-auto", middleware.NormalAuth(), reportHandler.GeneratePDFVendorStartFromLast)
	api.Get("/list-pdf", middleware.NormalAuth(), reportHandler.FindPDF)
	api.Get("/daily-vendor", middleware.NormalAuth(), reportHandler.GeneratePDFDailyReportVendor)
	api.Get("/daily-vendor-auto", middleware.NormalAuth(), reportHandler.GeneratePDFVendorDailyStartFromLast)
	api.Get("/generate-pdf-monthly", middleware.NormalAuth(), reportHandler.GeneratePDFVendorMonthly)
	api.Get("/generate-pdf-stock", middleware.NormalAuth(), reportHandler.GeneratePDFStock)

	// Option
	api.Get("/opt-check-item", optionHandler.OptCreateCheckItem)
	api.Get("/opt-stock", optionHandler.OptCreateStock)
	api.Get("/opt-cctv", optionHandler.OptCreateCctv)
	api.Get("/opt-computer", optionHandler.OptCreateComputer)
	api.Get("/opt-other", optionHandler.OptLocationDivision)
	api.Get("/opt-branch", optionHandler.OptBranch)
}

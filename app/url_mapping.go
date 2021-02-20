package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/muchlist/risa_restfull/middleware"
)

func mapUrls(app *fiber.App) {
	app.Use(logger.New())
	app.Use(middleware.LimitRequest())

	app.Static("/media", "./static/media")

	api := app.Group("/api/v1")
	api.Get("/ping", pingHandler.Ping)
	api.Post("/users", userHandler.Register)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)

	api.Get("/users", middleware.AuthMiddleware, userHandler.Find)
	api.Get("/profile", middleware.AuthMiddleware, userHandler.GetProfile)
	api.Post("/avatar", middleware.AuthMiddleware, userHandler.UploadImage)

	api.Post("/change-password", middleware.AuthMustFreshMiddleware, userHandler.ChangePassword)

	apiAuthAdmin := app.Group("/api/v1/admin")
	apiAuthAdmin.Use(middleware.AuthAdminMiddleware)
	apiAuthAdmin.Put("/users/:user_email", userHandler.Edit)
	apiAuthAdmin.Delete("/users/:user_email", userHandler.Delete)
	apiAuthAdmin.Get("/users/:user_email/reset-password", userHandler.ResetPassword)

}

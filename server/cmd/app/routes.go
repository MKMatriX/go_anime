package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

func (app *Application) routes() {
	app.server.Use(echoMiddleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	app.server.Use(echoMiddleware.Recover())       // recover panics as errors for proper error handling

	app.server.GET("/health", func(ctx *echo.Context) error {
		return ctx.String(http.StatusOK, "Echo is up")
	})

	app.server.POST("/api/v1/users", app.handler.UserRegister)
	app.server.POST("/api/v1/users/login", app.handler.UserLogin)

	app.server.GET("/api/v1/users/check", app.handler.GetAuthenticatedUser, app.appMiddleware.Authentication)
}

package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func (app *Application) routes() {
	app.server.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	app.server.Use(middleware.Recover())       // recover panics as errors for proper error handling

	app.server.GET("/health", func(ctx *echo.Context) error {
		return ctx.String(http.StatusOK, "Echo is up")
	})

	app.server.POST("/api/v1/users", app.handler.RegisterUser)
}

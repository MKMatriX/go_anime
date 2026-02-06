package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *Application) routes() {
	// app.server.Use(middleware.RequestLogger())

	app.server.GET("/health", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Echo is up")
	})

	app.server.GET("/api/v1/wallets/:WALLET_UUID", app.handler.GetBalance)
	app.server.POST("/api/v1/wallet/new", app.handler.CreateWallet)
	app.server.POST("/api/v1/wallet", app.handler.PerformOperation)
}

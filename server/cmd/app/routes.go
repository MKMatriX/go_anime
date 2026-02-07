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

	apiGroup := app.server.Group("/api")
	apiV1Group := apiGroup.Group("/v1")
	userGroup := apiV1Group.Group("/users")
	{
		userGroup.POST("", app.handler.UserRegister)
		userGroup.POST("/login", app.handler.UserLogin)
	}

	profileGroup := apiV1Group.Group("/profile", app.appMiddleware.Authentication)
	{
		profileGroup.GET("/check", app.handler.GetAuthenticatedUser)
	}

	animeGroup := apiV1Group.Group("/anime", app.appMiddleware.Authentication)
	{
		animeGroup.GET("", app.handler.AnimeList)
		animeGroup.POST("", app.handler.AnimeCreate)
		animeGroup.PUT("/:ID", app.handler.AnimeUpdate)
		animeGroup.DELETE("/:ID", app.handler.AnimeDelete)
	}

}

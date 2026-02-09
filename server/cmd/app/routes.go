package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func (app *Application) routes() {
	app.server.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	app.server.Use(middleware.Recover())       // recover panics as errors for proper error handling

	app.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		}, // В проде: []string{"https://your-frontend.com"}
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization, // ← обязательно для токенов/JWT
			// Добавь другие, если используешь (например, "X-Requested-With")
		},
		AllowCredentials: true, // ← если используешь куки, credentials: 'include' на фронте или Bearer-токены с auth
		// MaxAge:           300, // Опционально: кэшировать preflight-ответ на 5 минут (в секундах)
		// ExposeHeaders:    []string{"X-Custom-Header"}, // если фронту нужно читать кастомные заголовки из ответа
	}))

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
		profileGroup.GET("", app.handler.GetAuthenticatedUser)
	}

	animeGroup := apiV1Group.Group("/anime", app.appMiddleware.Authentication)
	{
		animeGroup.GET("", app.handler.AnimeList)
		animeGroup.GET("/:id", app.handler.AnimeItem)
		animeGroup.POST("", app.handler.AnimeCreate)
		animeGroup.PUT("/:id", app.handler.AnimeUpdate)
		animeGroup.DELETE("/:id", app.handler.AnimeDelete)
	}

}

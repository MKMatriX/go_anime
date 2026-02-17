package main

import (
	"fmt"
	"log"
	"log/slog"

	"go_anime/internal/handlers"
	"go_anime/internal/shared/init_db"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Application struct {
	server  *echo.Echo
	db      *gorm.DB
	handler *handlers.Handler
}

func main() {
	_ = godotenv.Load()

	dsn := init_db.GetDsn()
	db := init_db.InitDb(dsn)
	e := echo.New()

	h := handlers.NewHandler(db)

	app := Application{
		server:  e,
		db:      db,
		handler: h,
	}

	app.server.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	app.server.Use(middleware.Recover())       // recover panics as errors for proper error handling

	apiV1Group := app.server.Group("/api/v1")

	startMessage := fmt.Sprintf("Server starting on :%s", 8083)
	log.Println(startMessage)

	address := fmt.Sprintf(":%s", 8083)
	if err := e.Start(address); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

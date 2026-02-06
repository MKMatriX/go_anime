package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"go_anime/internal/handlers"
	"go_anime/internal/init_db"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
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
	init_db.MigrateDb(db)

	e := echo.New()

	h := handlers.NewHandler(db)

	app := Application{
		server:  e,
		db:      db,
		handler: h,
	}

	app.routes()

	startMessage := fmt.Sprintf("Server starting on :%s", os.Getenv("INTERNAL_APP_PORT"))
	log.Println(startMessage)

	address := fmt.Sprintf(":%s", os.Getenv("INTERNAL_APP_PORT"))
	if err := e.Start(address); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

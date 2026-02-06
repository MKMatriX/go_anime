package main

import (
	"fmt"
	"log"
	"os"

	"go_anime/internal/dbSetup"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Application struct {
	server *echo.Echo
	db     *gorm.DB
}

func main() {
	_ = godotenv.Load()

	dsn := dbSetup.GetDsn()
	db := dbSetup.InitDb(dsn)
	dbSetup.MigrateDb(dsn)

	e := echo.New()

	app := Application{
		server: e,
		db:     db,
	}

	app.routes()

	startMessage := fmt.Sprintf("Server starting on :%s", os.Getenv("INTERNAL_APP_PORT"))
	log.Println(startMessage)

	address := fmt.Sprintf(":%s", os.Getenv("INTERNAL_APP_PORT"))
	e.Logger.Fatal(e.Start(address))
}

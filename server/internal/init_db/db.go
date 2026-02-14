package init_db

import (
	"fmt"
	"go_anime/internal/models"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func InitDb(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	return db
}

func MigrateDb(db *gorm.DB) {
	err := db.AutoMigrate(&models.UserModel{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.AnimeModel{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.AnimeEpisodeModel{})

	if err != nil {
		panic(err)
	}
}

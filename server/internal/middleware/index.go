package middleware

import "gorm.io/gorm"

type AppMiddleware struct {
	DB *gorm.DB
}

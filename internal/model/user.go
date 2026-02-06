package model

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Login    *string `gorm:"type:varchar(255);uniqueIndex" json:"login"`
	Password string  `gorm:"type:varchar(255)" json:"-"`
}

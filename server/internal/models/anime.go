package models

import "gorm.io/gorm"

type AnimeModel struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (receiver AnimeModel) TableName() string {
	return "anime"
}

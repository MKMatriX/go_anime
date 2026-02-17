package models

type UserModel struct {
	BaseModel
	Login            string `gorm:"type:varchar(255);uniqueIndex;not null" json:"login"`
	Password         string `gorm:"type:varchar(255);not null" json:"-"`
	RefreshTokenHash string `gorm:"type:varchar(255);index" json:"-"`
}

func (receiver UserModel) TableName() string {
	return "users"
}

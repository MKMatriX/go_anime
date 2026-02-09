package models

type AnimeModel struct {
	BaseModel
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
	AnilistInfo string `json:"anilistInfo"`
}

func (receiver AnimeModel) TableName() string {
	return "anime"
}

package models

type AnimeModel struct {
	BaseModel
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
	AnilistInfo string `json:"anilistInfo"`
	ShikiInfo   string `json:"shikiInfo"`
}

func (receiver AnimeModel) TableName() string {
	return "anime"
}

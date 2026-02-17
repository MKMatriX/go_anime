package models

type AnimeModel struct {
	BaseModel
	Name        string              `json:"name" validate:"required,min=1"`
	Description string              `json:"description"`
	AnilistInfo string              `json:"anilistInfo"`
	ShikiInfo   string              `json:"shikiInfo"`
	AniDBId     uint                `json:"aniDBId"`
	Episodes    []AnimeEpisodeModel `json:"episodes" gorm:"foreignKey:AnimeID"`
}

func (receiver AnimeModel) TableName() string {
	return "anime"
}

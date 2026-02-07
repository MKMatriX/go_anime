package models

type AnimeModel struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (receiver AnimeModel) TableName() string {
	return "anime"
}

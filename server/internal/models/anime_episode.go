package models

type AnimeEpisodeModel struct {
	BaseModel

	Anime   AnimeModel
	AnimeID uint `json:"animeId gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	EpisodeNumber int    `json:"episodeNumber"`
	Name          string `json:"name"`
	Translator    string `json:"translator"`
	Width         string `json:"width"`
	TorrentUrl    string `json:"torrentUrl"`
	MagnetUrl     string `json:"magnetUrl"`
	LocalUrl      string `json:"localUrl"`
}

func (receiver AnimeEpisodeModel) TableName() string {
	return "anime_episode"
}

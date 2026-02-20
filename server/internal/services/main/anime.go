package services

import (
	"fmt"
	"go_anime/internal/requests"
	"go_anime/internal/shared/models"

	"gorm.io/gorm"
)

type AnimeSevice struct {
	db *gorm.DB
}

func NewAnimeService(db *gorm.DB) *AnimeSevice {
	return &AnimeSevice{db: db}
}

func (s *AnimeSevice) List() []*models.AnimeModel {
	var anime []*models.AnimeModel
	s.db.Find(&anime)
	return anime
}

func (s *AnimeSevice) GetById(id uint) (*models.AnimeModel, error) {
	var anime *models.AnimeModel
	result := s.db.Preload("Episodes").Where("id = ?", id).First(&anime)
	if result.Error != nil {
		return nil, result.Error
	}
	return anime, nil
}

func (s *AnimeSevice) Create(request *requests.AnimeCreateRequest) (*models.AnimeModel, error) {
	anime := models.AnimeModel{
		Name:        request.Name,
		Description: request.Description,
	}
	result := s.db.Create(&anime)

	pipe := func() {
		s.getAniDBId(&anime)
		// we need anidbId for episode search
		s.GetEpisodes(&anime)
	}
	go pipe()
	go s.getAnilibInfo(&anime) // yey my first go routine in this project
	go s.getShikiInfo(&anime)

	if result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

func (s *AnimeSevice) Update(id uint, request *requests.AnimeCreateRequest) (*models.AnimeModel, error) {
	var anime *models.AnimeModel
	result := s.db.Where("id = ?", id).First(&anime)
	if result.Error != nil {
		return nil, result.Error
	}

	anime.Name = request.Name
	anime.Description = request.Description
	result = s.db.Save(anime)
	if result.Error != nil {
		return nil, result.Error
	}

	return anime, nil
}

func (s *AnimeSevice) Delete(id uint) error {
	result := s.db.Delete(&models.AnimeModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("anime not found")
	}
	return nil
}

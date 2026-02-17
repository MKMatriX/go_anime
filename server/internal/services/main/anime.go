package services

import (
	"context"
	"fmt"
	"go_anime/internal/requests"
	"go_anime/internal/shared/models"
	"go_anime/internal/shared/proto/anidb"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		// s.GetEpisodes(&anime)
	}
	go pipe()
	// go s.getAnilibInfo(&anime) // yey my first go routine in this project
	// go s.getShikiInfo(&anime)

	if result.Error != nil {
		return nil, result.Error
	}
	return &anime, nil
}

// func (s *AnimeSevice) GetEpisodes(anime *models.AnimeModel) ([]*models.AnimeEpisodeModel, error) {
// 	var dbEpisodes []*models.AnimeEpisodeModel

// 	toshoItems, err := GetToshoEpisodes(anime.AniDBId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, toshoItem := range toshoItems {
// 		episode, ok := ParseToshoItemToEpisode(toshoItem, anime.ID)
// 		if ok {
// 			dbEpisodes = append(dbEpisodes, &episode)
// 		} else {
// 			fmt.Println("Failed to parse Tosho item ", toshoItem)
// 		}
// 	}

// 	result := s.db.Create(dbEpisodes)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return dbEpisodes, nil
// }

// func (s *AnimeSevice) getAnilibInfo(anime *models.AnimeModel) {
// 	info, err := GetAnilistAnimeInfo(anime.Name)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return
// 	}

// 	json, err := json.Marshal(info)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return
// 	}

// 	anime.AnilistInfo = string(json)
// 	result := s.db.Save(anime)
// 	if result.Error != nil {
// 		slog.Error(result.Error.Error())
// 	}
// }

// func (s *AnimeSevice) getShikiInfo(anime *models.AnimeModel) {
// 	info, err := GetShikiAnimeInfo(anime.Name)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return
// 	}

// 	json, err := json.Marshal(info)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return
// 	}

// 	anime.ShikiInfo = string(json)
// 	result := s.db.Save(anime)
// 	if result.Error != nil {
// 		slog.Error(result.Error.Error())
// 	}
// }

func (s *AnimeSevice) getAniDBId(anime *models.AnimeModel) {
	conn, err := grpc.NewClient("anidb:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to dial anidb: " + err.Error())
		return
	}
	defer conn.Close()

	client := anidb.NewAniDBServiceClient(conn)
	resp, err := client.GetAniDBId(
		context.Background(),
		&anidb.GetAniDBIdRequest{AnimeName: anime.Name},
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if resp.Error != "" {
		slog.Error(resp.Error)
		return
	}

	anime.AniDBId = uint(resp.Id)
	result := s.db.Save(anime)
	if result.Error != nil {
		slog.Error(result.Error.Error())
	}
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

package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_anime/internal/requests"
	"go_anime/internal/shared/models"
	"go_anime/internal/shared/proto/anidb"
	"go_anime/internal/shared/proto/anilist"
	"go_anime/internal/shared/proto/animetosho"
	"log/slog"
	"os"

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

func (s *AnimeSevice) GetEpisodes(anime *models.AnimeModel) ([]*models.AnimeEpisodeModel, error) {
	conn, err := grpc.NewClient(
		"animetosho:"+os.Getenv("ANIMETOSHO_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("failed to dial animetosho: " + err.Error())
		return nil, err
	}
	defer conn.Close()

	client := animetosho.NewAnimeToshoServiceClient(conn)
	resp, err := client.GetParsedEpisodes(
		context.Background(),
		&animetosho.GetParsedEpisodesRequest{
			AnidbId: uint32(anime.AniDBId),
			AnimeId: uint32(anime.ID),
		},
	)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	if resp.ErrorMessage != "" {
		slog.Error(resp.ErrorMessage)
		return nil, errors.New(resp.ErrorMessage)
	}
	if resp.Episodes == nil || len(resp.Episodes) == 0 {
		slog.Error("No episodes")
		return nil, errors.New("No episodes")
	}

	var modelEpisodes []*models.AnimeEpisodeModel
	// Сохраняем эпизоды в БД
	for _, ep := range resp.Episodes {
		model := models.AnimeEpisodeModel{
			AnimeID:       uint(ep.AnimeId),
			EpisodeNumber: int(ep.EpisodeNumber),
			Name:          ep.Name,
			Translator:    ep.Translator,
			Width:         ep.Width,
			TorrentUrl:    ep.TorrentUrl,
			MagnetUrl:     ep.MagnetUrl,
			LocalUrl:      ep.LocalUrl,
		}
		modelEpisodes = append(modelEpisodes, &model)
	}

	// result := s.db.Create(modelEpisodes)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	return modelEpisodes, nil
}

func (s *AnimeSevice) getAnilibInfo(anime *models.AnimeModel) {
	conn, err := grpc.NewClient(
		"anilist:"+os.Getenv("ANILIST_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("failed to dial anilist: " + err.Error())
		return
	}
	defer conn.Close()

	client := anilist.NewAniListServiceClient(conn)
	resp, err := client.GetAnimeInfo(
		context.Background(),
		&anilist.GetAnimeInfoRequest{Title: anime.Name},
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if resp.GetError() != nil {
		slog.Error(resp.GetError().Message)
		return
	}

	json, err := json.Marshal(resp.Result)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	anime.AnilistInfo = string(json)
	result := s.db.Save(anime)
	if result.Error != nil {
		slog.Error(result.Error.Error())
	}
}

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
	conn, err := grpc.NewClient(
		"anidb:"+os.Getenv("ANIDB_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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

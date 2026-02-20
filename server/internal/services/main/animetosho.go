package services

import (
	"context"
	"errors"
	"go_anime/internal/shared/models"
	"go_anime/internal/shared/proto/animetosho"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
	if !resp.HasEpisodes {
		slog.Error("No episodes")
		return nil, errors.New("No episodes")
	}

	var modelEpisodes []*models.AnimeEpisodeModel
	// Сохраняем эпизоды в БД
	for _, ep := range resp.GetEpisodes() {
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

	result := s.db.Create(modelEpisodes)
	if result.Error != nil {
		return nil, result.Error
	}

	return modelEpisodes, nil
}

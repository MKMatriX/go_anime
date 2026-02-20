package services

import (
	"context"
	"encoding/json"
	"go_anime/internal/shared/models"
	"go_anime/internal/shared/proto/shikimori"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *AnimeSevice) getShikiInfo(anime *models.AnimeModel) {
	conn, err := grpc.NewClient(
		"shikimori:"+os.Getenv("SHIKIMORI_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("failed to dial shikimori: " + err.Error())
		return
	}
	defer conn.Close()

	client := shikimori.NewShikimoriServiceClient(conn)
	resp, err := client.GetAnimeInfo(
		context.Background(),
		&shikimori.GetAnimeInfoRequest{Search: anime.Name},
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if resp.GetError() != nil {
		slog.Error(resp.GetError().Message)
		return
	}

	json, err := json.Marshal(resp.GetAnime())
	if err != nil {
		slog.Error(err.Error())
		return
	}

	anime.ShikiInfo = string(json)
	result := s.db.Save(anime)
	if result.Error != nil {
		slog.Error(result.Error.Error())
	}
}

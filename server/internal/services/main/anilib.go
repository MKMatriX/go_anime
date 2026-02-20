package services

import (
	"context"
	"encoding/json"
	"go_anime/internal/shared/models"
	"go_anime/internal/shared/proto/anilist"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

	json, err := json.Marshal(resp.GetAnime())
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

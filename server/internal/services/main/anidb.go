package services

import (
	"context"
	"go_anime/internal/shared/models"
	"go_anime/internal/shared/proto/anidb"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

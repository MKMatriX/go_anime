package services

import (
	"context"
	"errors"
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

type TitleEntry struct {
	AID   int
	Type  string // "main", "official", "syn", "short", "kana" и т.д.
	Lang  string
	Title string
}

func (s *AnimeSevice) AutocompleteSearch(query string) ([]*TitleEntry, error) {
	conn, err := grpc.NewClient(
		"anidb:"+os.Getenv("ANIDB_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("failed to dial anidb: " + err.Error())
		return nil, err
	}
	defer conn.Close()

	client := anidb.NewAniDBServiceClient(conn)
	resp, err := client.AutocompleteSearch(
		context.Background(),
		&anidb.GetAniDBIdRequest{AnimeName: query},
	)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	if resp.Error != "" {
		slog.Error(resp.Error)
		return nil, errors.New(resp.Error)
	}

	protoTitles := resp.GetTitles()
	var titles []*TitleEntry
	for _, pt := range protoTitles {
		titles = append(titles, &TitleEntry{
			AID:   int(pt.AnimeId),
			Type:  pt.Type,
			Lang:  pt.Lang,
			Title: pt.Title,
		})
	}

	return titles, nil
}

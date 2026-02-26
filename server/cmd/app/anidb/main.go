package main

import (
	"context"
	"log"
	"net"
	"os"

	anidbService "go_anime/internal/services/anidb" // Ваша логика
	services "go_anime/internal/services/anidb"
	"go_anime/internal/shared/proto/anidb" // Сгенерированный proto

	"google.golang.org/grpc"
)

type server struct {
	anidb.UnimplementedAniDBServiceServer
}

func (s *server) GetAniDBId(ctx context.Context, req *anidb.GetAniDBIdRequest) (*anidb.GetAniDBIdResponse, error) {
	id, err := anidbService.GetAniDBId(req.AnimeName) // Ваша оригинальная функция
	if err != nil {
		return &anidb.GetAniDBIdResponse{Error: err.Error()}, nil
	}
	return &anidb.GetAniDBIdResponse{Id: int32(id)}, nil
}

func (s *server) AutocompleteSearch(ctx context.Context, req *anidb.GetAniDBIdRequest) (*anidb.AnimeTitles, error) {
	titles, err := anidbService.AutocompleteSearch(req.AnimeName) // Ваша оригинальная функция
	if err != nil {
		return &anidb.AnimeTitles{Error: err.Error()}, nil
	}
	protoTitles := titlesToProtoTitles(titles)
	return &anidb.AnimeTitles{Titles: protoTitles}, nil
}

func (s *server) GetOtherNames(ctx context.Context, req *anidb.AnimeIdRequest) (*anidb.AnimeTitles, error) {
	titles, err := anidbService.GetOtherNames(int(req.AnimeId))
	if err != nil {
		return &anidb.AnimeTitles{Error: err.Error()}, nil
	}
	protoTitles := titlesToProtoTitles(titles)
	return &anidb.AnimeTitles{Titles: protoTitles}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("ANIDB_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	anidb.RegisterAniDBServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func titleToProtoTitle(title services.TitleEntry) *anidb.AnimeTitle {
	return &anidb.AnimeTitle{
		AnimeId: int32(title.AID),
		Type:    title.Type,
		Title:   title.Title,
		Lang:    title.Lang,
	}
}

func titlesToProtoTitles(titles []services.TitleEntry) []*anidb.AnimeTitle {
	var protoTitles []*anidb.AnimeTitle
	for _, t := range titles {
		protoTitles = append(protoTitles, titleToProtoTitle(t))
	}
	return protoTitles
}

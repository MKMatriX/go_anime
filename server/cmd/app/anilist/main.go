package main

import (
	"context"
	"log"
	"net"
	"os"

	// Ваша логика
	services "go_anime/internal/services/anilist"
	"go_anime/internal/shared/proto/anilist"

	"google.golang.org/grpc"
)

type anilistServer struct {
	anilist.UnimplementedAniListServiceServer
}

func (s *anilistServer) GetAnimeInfo(ctx context.Context, req *anilist.GetAnimeInfoRequest) (*anilist.GetAnimeInfoResponse, error) {
	info, err := services.GetAnilistAnimeInfo(req.Title) // ваша текущая функция
	if err != nil {
		return &anilist.GetAnimeInfoResponse{
			Result: &anilist.GetAnimeInfoResponse_Error{
				Error: &anilist.Error{
					Message: err.Error(),
					Code:    5, // например 5 = NOT_FOUND
				},
			},
		}, nil
	}

	// маппинг из вашей структуры в proto-сообщение
	protoInfo := &anilist.AnimeInfo{
		Id: int32(info.ID),
		Title: &anilist.Title{
			Romaji:  info.Title.Romaji,
			English: info.Title.English,
			Native:  info.Title.Native,
		},
		Description: info.Description,
		Format:      info.Format,
		Status:      info.Status,

		// Просто присваиваем — если было nil → будет 0
		Episodes:     int32OrZero(info.Episodes),
		Duration:     int32OrZero(info.Duration),
		AverageScore: int32OrZero(info.AverageScore),
		Genres:       info.Genres,
		Season:       info.Season,
		SeasonYear:   int32OrZero(info.SeasonYear),

		StartDate: &anilist.Date{
			Year:  int32OrZero(info.StartDate.Year),
			Month: int32OrZero(info.StartDate.Month),
			Day:   int32OrZero(info.StartDate.Day),
		},
		CoverImage: &anilist.Image{
			Large: info.CoverImage.Large,
		},
	}

	return &anilist.GetAnimeInfoResponse{
		Result: &anilist.GetAnimeInfoResponse_Anime{
			Anime: protoInfo,
		},
	}, nil
}

// Вспомогательные конвертеры
func int32OrZero(v *int) int32 {
	if v == nil {
		return 0
	}
	return int32(*v)
}

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("ANILIST_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	anilist.RegisterAniListServiceServer(s, &anilistServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

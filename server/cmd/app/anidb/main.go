package main

import (
	"context"
	"log"
	"net"

	anidbService "go_anime/internal/services/anidb" // Ваша логика
	"go_anime/internal/shared/proto/anidb"          // Сгенерированный proto

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

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	anidb.RegisterAniDBServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

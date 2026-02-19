package main

import (
	"context"
	"log"
	"net"
	"os"

	animetoshoService "go_anime/internal/services/animetosho"
	"go_anime/internal/shared/proto/animetosho"
	pb "go_anime/internal/shared/proto/animetosho"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAnimeToshoServiceServer
}

func (s *server) GetParsedEpisodes(ctx context.Context, req *pb.GetParsedEpisodesRequest) (*pb.GetParsedEpisodesResponse, error) {
	log.Printf("AnimeTosho request: AniDB %d → Anime %d", req.AnidbId, req.AnimeId)

	episodes, err := animetoshoService.GetParsedToshoEpisodes(
		uint(req.AnidbId),
		uint(req.AnimeId),
	)

	resp := &pb.GetParsedEpisodesResponse{
		HasEpisodes: len(episodes) > 0,
	}

	if err != nil {
		resp.ErrorMessage = err.Error()
		return resp, nil
	}

	for _, model := range episodes {
		protoEp := &pb.AnimeEpisode{
			AnimeId:       uint32(model.AnimeID),
			EpisodeNumber: int32(model.EpisodeNumber),
			Name:          model.Name,
			Translator:    model.Translator,
			Width:         model.Width,
			TorrentUrl:    model.TorrentUrl,
			MagnetUrl:     model.MagnetUrl,
			LocalUrl:      model.LocalUrl,
		}
		resp.Episodes = append(resp.Episodes, protoEp)
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("ANIMETOSHO_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	animetosho.RegisterAnimeToshoServiceServer(s, &server{})
	log.Printf("AnimeTosho gRPC server listening on :" + os.Getenv("ANIMETOSHO_PORT"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

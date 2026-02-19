package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	shikimori "go_anime/internal/services/shikimori" // ваш текущий пакет с GetShikiAnimeInfo
	pb "go_anime/internal/shared/proto/shikimori"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedShikimoriServiceServer
}

func (s *server) GetAnimeInfo(ctx context.Context, req *pb.GetAnimeInfoRequest) (*pb.GetAnimeInfoResponse, error) {
	log.Printf("Shikimori request: search = %q", req.Search)

	anime, err := shikimori.GetShikiAnimeInfo(req.Search)
	if err != nil {
		slog.Error(err.Error())
		return &pb.GetAnimeInfoResponse{
			Result: &pb.GetAnimeInfoResponse_Error{
				Error: &pb.Error{Message: err.Error()},
			},
		}, nil
	}

	protoAnime := &pb.ShikiAnime{
		Id:                anime.ID,
		Name:              anime.Name,
		Russian:           anime.Russian,
		English:           anime.English,
		Description:       anime.Description,
		DescriptionHtml:   anime.DescriptionHtml,
		DescriptionSource: anime.DescriptionSource,
		Score:             float32(anime.Score), // float64 → float32 (proto использует float32)
		Kind:              anime.Kind,
		Status:            anime.Status,
		Episodes:          int32(anime.Episodes),
		Poster: &pb.Poster{
			Id:          anime.Poster.ID,
			OriginalUrl: anime.Poster.OriginalUrl,
			MainUrl:     anime.Poster.MainUrl,
		},
	}

	for _, role := range anime.PersonRoles {
		pRole := &pb.PersonRole{
			Id:      role.ID,
			RolesRu: role.RolesRu,
			RolesEn: role.RolesEn,
			Person: &pb.Person{
				Id:   role.Person.ID,
				Name: role.Person.Name,
			},
		}
		if role.Person.Poster != nil {
			pRole.Person.Poster = &pb.Poster{
				Id:          role.Person.Poster.ID,
				OriginalUrl: role.Person.Poster.OriginalUrl,
				MainUrl:     role.Person.Poster.MainUrl,
			}
		}
		protoAnime.PersonRoles = append(protoAnime.PersonRoles, pRole)
	}

	for _, link := range anime.ExternalLinks {
		protoLink := &pb.ExternalLink{
			Kind: link.Kind,
			Url:  link.Url,
		}
		if link.ID != nil {
			protoLink.Id = *link.ID
		}
		protoAnime.ExternalLinks = append(protoAnime.ExternalLinks, protoLink)
	}

	return &pb.GetAnimeInfoResponse{
		Result: &pb.GetAnimeInfoResponse_Anime{
			Anime: protoAnime,
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("SHIKIMORI_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterShikimoriServiceServer(s, &server{})

	log.Println("Shikimori gRPC server listening on :" + os.Getenv("SHIKIMORI_PORT"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

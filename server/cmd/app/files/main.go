package main

import (
	"context"
	services "go_anime/internal/services/files"
	"go_anime/internal/shared/proto/files"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type fileServiceServer struct {
	files.UnimplementedFileServiceServer
}

func (*fileServiceServer) DownloadFile(ctx context.Context, req *files.DownloadFileRequest) (*files.DownloadFileResponse, error) {
	path, err := services.DownloadFile(req.Url)
	if err != nil {
		return nil, err
	}

	protoStruct := &files.DownloadFileResponse{
		Path: path,
	}

	return protoStruct, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("FILES_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	files.RegisterFileServiceServer(s, &fileServiceServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

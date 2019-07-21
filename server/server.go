package server

import (
	"context"
	"fmt"
	"net"

	"satelit-project/satelit-scraper/proto/scraper"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type scraperServiceServer struct{}

func Serve(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	scraper.RegisterScraperServiceServer(grpcServer, &scraperServiceServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (s *scraperServiceServer) StartScraping(ctx context.Context, intent *scraper.ScrapeIntent) (*empty.Empty, error) {
	err := RunScraper(ctx, intent)
	if err != nil {
		return nil, err
	}

	stub := &empty.Empty{}
	return stub, nil
}

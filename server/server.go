package server

import (
	"context"
	"fmt"
	"net"

	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"

	"google.golang.org/grpc"
)

type scraperServiceServer struct{}

func Serve(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	scraping.RegisterScraperServiceServer(grpcServer, &scraperServiceServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (s *scraperServiceServer) StartScraping(ctx context.Context, intent *scraping.ScrapeIntent) (*scraping.ScrapeIntentResult, error) {
	mayContinue, err := RunScraper(ctx, intent)
	if err != nil {
		return nil, err
	}

	result := &scraping.ScrapeIntentResult{
		Id:          intent.Id,
		MayContinue: mayContinue,
	}

	return result, nil
}

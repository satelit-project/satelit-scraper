package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
	"shitty.moe/satelit-project/satelit-scraper/proxy"
	"shitty.moe/satelit-project/satelit-scraper/proxy/provider"
)

// TODO: move to config
const _tasksAddr string = "127.0.0.1:10602"

type scraperServiceServer struct {
	tasksClient *grpc.ClientConn
	proxies     *proxy.Fetcher
}

func Serve(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(_tasksAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	robinProvider := provider.NewRoundRobin([]proxy.Provider{
		provider.NewPLD(),
		provider.NewPSC(),
	})

	grpc := grpc.NewServer()
	server := scraperServiceServer{
		tasksClient: conn,
		proxies:     proxy.NewFetcher(robinProvider, ProxiesLimit, proxy.HTTP),
	}

	scraping.RegisterScraperServiceServer(grpc, &server)

	logging.DefaultLogger().Infof("Start listening on port %d", port)
	err = grpc.Serve(lis)
	if err != nil {
		return err
	}

	return conn.Close()
}

func (s *scraperServiceServer) StartScraping(ctx context.Context, intent *scraping.ScrapeIntent) (*scraping.ScrapeIntentResult, error) {
	runner := NewRunner(s.tasksClient, s.proxies, logging.DefaultLogger())
	mayContinue, err := runner.Run(ctx, intent)
	if err != nil {
		return nil, err
	}

	result := &scraping.ScrapeIntentResult{
		Id:          intent.Id,
		MayContinue: mayContinue,
	}

	return result, nil
}

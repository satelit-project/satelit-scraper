package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"shitty.moe/satelit-project/satelit-scraper/config"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
)

// GRPC service for anime scraping.
type ScrapingService struct {
	inner *grpc.Server
	cfg   config.Config
	log   *logging.Logger
}

// Creates new scraping service instance.
func New(cfg config.Config, log *logging.Logger) *ScrapingService {
	return &ScrapingService{
		inner: grpc.NewServer(),
		cfg:   cfg,
		log:   log,
	}
}

// Synchronously runs the service.
func (s *ScrapingService) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Serving.Port))
	if err != nil {
		return err
	}

	scraping.RegisterScraperServiceServer(s.inner, s)
	if err = s.inner.Serve(lis); err != nil {
		return err
	}

	return nil
}

// Gracefully shuts down the service.
func (s *ScrapingService) Shutdown() {
	s.log.Infof("trying to gracefully shutdown the server")
	s.inner.GracefulStop()
}

// Starts scraping process for given intent.
func (s *ScrapingService) StartScraping(ctx context.Context, intent *scraping.ScrapeIntent) (*scraping.ScrapeIntentResult, error) {
	runner := NewRunner(s.cfg.Scraping, s.cfg.AniDB, s.log)
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

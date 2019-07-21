package spider

import (
	"context"

	"satelit-project/satelit-scraper/proto/scraper"

	"github.com/golang/protobuf/ptypes/empty"
)

type scraperServiceServer struct{}

func (s *scraperServiceServer) StartScraping(ctx context.Context, intent *scraper.ScrapeIntent) (*empty.Empty, error) {
	err := RunScraper(ctx, intent)
	if err != nil {
		return nil, err
	}

	stub := &empty.Empty{}
	return stub, nil
}

package server

import (
	"context"
	"time"

	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
)

// gRPC transport for scraping progress reporter.
type grpcTransport struct {
	client scraping.ScraperTasksServiceClient
	timeout time.Duration
}

func (g grpcTransport) Yield(ty *scraping.TaskYield) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	_, err := g.client.YieldResult(ctx, ty)
	cancel()
	return err
}

func (g grpcTransport) Finish(tf *scraping.TaskFinish) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	_, err := g.client.CompleteTask(ctx, tf)
	cancel()
	return err
}

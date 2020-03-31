package server

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"shitty.moe/satelit-project/satelit-scraper/config"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
	"shitty.moe/satelit-project/satelit-scraper/proxy"
	"shitty.moe/satelit-project/satelit-scraper/proxy/provider"
	"shitty.moe/satelit-project/satelit-scraper/spider"
	"shitty.moe/satelit-project/satelit-scraper/spider/anidb"
)

// Number of proxies to fetch for scraping.
const ProxiesLimit int = 8

// Number of jobs to request per spider run.
const TitlesPerRun int32 = 8

// Runner for AniDB spider.
type AniDBRunner struct {
	cfg      *config.Scraping
	anidbCfg *config.AniDB
	cache    spider.Cache
	log      *logging.Logger
}

// Creates new runner instance.
func NewRunner(cfg *config.Scraping, anidbCfg *config.AniDB, cache spider.Cache, log *logging.Logger) AniDBRunner {
	return AniDBRunner{
		cfg:      cfg,
		anidbCfg: anidbCfg,
		cache:    cache,
		log:      log,
	}
}

// Runs AniDB spider for provided intent. Returns true if there's more data to scrape.
func (r AniDBRunner) Run(context context.Context, intent *scraping.ScrapeIntent) (bool, error) {
	log := r.log.With("anidb-intent", intent.Id)
	log.Infof("received scraping intent: %s", intent.Id)

	conn, err := grpc.Dial(r.cfg.TaskAddress, grpc.WithInsecure())
	if err != nil {
		return false, nil
	}

	client := scraping.NewScraperTasksServiceClient(conn)
	defer conn.Close()

	cmd := scraping.TaskCreate{
		Limit:  TitlesPerRun,
		Source: intent.Source,
	}

	task, err := client.CreateTask(context, &cmd)
	if err != nil {
		log.Errorf("failed to create scraping task: %v", err)
		return false, err
	}

	if len(task.Jobs) == 0 {
		log.Infof("task is empty: %v", task.Id)
		return false, nil
	}

	startAniDBScraping(spiderContext{
		intent: intent,
		task:   task,
		client: client,
		cache:  r.cache,
		cfg:    r.anidbCfg,
		log:    log,
	})

	return true, nil
}

// Context for running AniDB spiders.
type spiderContext struct {
	intent *scraping.ScrapeIntent
	task   *scraping.Task
	client scraping.ScraperTasksServiceClient
	cache  spider.Cache
	cfg    *config.AniDB
	log    *logging.Logger
}

// Starts AniDB spider with data from provided context.
func startAniDBScraping(ctx spiderContext) {
	log := ctx.log
	providers := provider.NewRoundRobin([]proxy.Provider{
		provider.NewPLD(),
		provider.NewPSC(),
	})

	fetcher := proxy.NewFetcher(providers, proxyLimit(ctx.task), proxy.HTTP, log)
	proxies := fetcher.Fetch()

	tr := grpcTransport{ctx.client, time.Duration(ctx.cfg.Timeout) * time.Second}
	reporter := spider.TaskReporter{Task: ctx.task, Transport: tr}
	spdr := anidb.NewSpider(&reporter, ctx.cache, ctx.cfg, log)

	if len(proxies) == 0 {
		log.Errorf("no proxies fetched, skipping")
		if err := reporter.Finish(); err != nil {
			log.Errorf("failed to report task finish: %v", err)
		}
		return
	}

	spdr.SetProxies(proxies)
	spdr.Run()
}

// Returns desired number of proxies to fetch for the task.
func proxyLimit(t *scraping.Task) int {
	if len(t.Jobs) < ProxiesLimit {
		return len(t.Jobs)
	}

	return ProxiesLimit
}

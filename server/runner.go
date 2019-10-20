package server

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"shitty.moe/satelit-project/satelit-scraper/proto/data"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
	"shitty.moe/satelit-project/satelit-scraper/proxy"
	"shitty.moe/satelit-project/satelit-scraper/spider"
	"shitty.moe/satelit-project/satelit-scraper/spider/anidb"
)

const RunnersLimit int = 16
const ProxiesLimit int = 8
const TitlesPerRun int32 = 8

var limit = make(chan bool, RunnersLimit)

type SpiderRunner struct {
	conn         *grpc.ClientConn
	proxyFetcher *proxy.Fetcher
	log          *zap.SugaredLogger
}

func NewRunner(conn *grpc.ClientConn, proxy *proxy.Fetcher, log *zap.SugaredLogger) SpiderRunner {
	return SpiderRunner{
		conn:         conn,
		proxyFetcher: proxy,
		log:          log,
	}
}

func (s SpiderRunner) Run(context context.Context, intent *scraping.ScrapeIntent) (bool, error) {
	log := s.log.With("scraping-intent", intent.Id)
	log.Info("received scraping intent")

	limit <- true
	client := scraping.NewScraperTasksServiceClient(s.conn)

	cmd := scraping.TaskCreate{
		Limit:  TitlesPerRun,
		Source: intent.Source,
	}

	task, err := client.CreateTask(context, &cmd)
	if err != nil {
		<-limit
		log.Errorf("failed to create scraping task: %v", err)
		return false, err
	}

	if len(task.Jobs) == 0 {
		<-limit
		log.Infof("task is empty: %v", task.Id)
		return false, nil
	}

	go func() {
		runScraper(scrapeContext{
			intent: intent,
			task:   task,
			client: client,
			log:    log.With("task", task.Id),
		})
		<-limit
	}()

	return true, nil
}

type grpcTransport struct {
	client scraping.ScraperTasksServiceClient
}

func (g grpcTransport) Yield(ty *scraping.TaskYield) error {
	_, err := g.client.YieldResult(context.Background(), ty)
	return err
}

func (g grpcTransport) Finish(tf *scraping.TaskFinish) error {
	_, err := g.client.CompleteTask(context.Background(), tf)
	return err
}

type scrapeContext struct {
	intent  *scraping.ScrapeIntent
	task    *scraping.Task
	client  scraping.ScraperTasksServiceClient
	proxies *proxy.Fetcher
	log     *zap.SugaredLogger
}

func runScraper(ctx scrapeContext) {
	switch ctx.intent.Source {
	case data.Source_ANIDB:
		startAniDBScraping(ctx)
	default:
		ctx.log.Errorf("scraper for source is not implemented: %v", ctx.intent.Source)
	}
}

func startAniDBScraping(ctx scrapeContext) {
	proxies := ctx.proxies.Fetch()

	tr := grpcTransport{client: ctx.client}
	reporter := spider.NewTaskReporter(ctx.task, tr)
	spdr := anidb.NewSpider(ctx.task, reporter)

	if len(proxies) == 0 {
		ctx.log.Error("no proxies fetched, skipping")
		reporter.Finish()
		return
	}

	spdr.SetProxies(proxies)
	spdr.Run()
}

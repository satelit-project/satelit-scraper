package spider

import (
	"context"
	"fmt"

	"satelit-project/satelit-scraper/proto/scraper"
	"satelit-project/satelit-scraper/proxy"
	"satelit-project/satelit-scraper/proxy/provider"
	"satelit-project/satelit-scraper/spider/anidb"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const RunnersLimit int = 16
const ProxiesLimit int = 8

var runner *spiderRunner
var limit = make(chan bool, RunnersLimit)

type spiderRunner struct {
	conn         *grpc.ClientConn
	proxyFetcher *proxy.Fetcher
	log          *logrus.Entry
}

type grpcTransport struct {
	client scraper.ScraperTasksServiceClient
}

func (g grpcTransport) Yield(ty *scraper.TaskYield) error {
	_, err := g.client.YieldResult(context.Background(), ty)
	return err
}

func (g grpcTransport) Finish(tf *scraper.TaskFinish) error {
	_, err := g.client.CompleteTask(context.Background(), tf)
	return err
}

func Init(taskServerAddr string) {
	Deinit()

	conn, err := grpc.Dial(taskServerAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to initiate connection to %s: %v\n", taskServerAddr, err))
	}

	fetcher := proxy.NewFetcher(provider.NewPLD(), ProxiesLimit, proxy.HTTP)
	log := logrus.NewEntry(logrus.StandardLogger())

	runner = &spiderRunner{
		conn:         conn,
		proxyFetcher: fetcher,
		log:          log,
	}
}

func Deinit() {
	runner := runner
	if runner == nil {
		return
	}

	err := runner.conn.Close()
	if err != nil {
		runner.log.Warnf("failed to close client grpc connection: %v", err)
	}
}

func RunScraper(context context.Context, intent *scraper.ScrapeIntent) error {
	if runner == nil {
		panic("spider runner is not initialized")
	}

	log := runner.log.WithField("scraping-intent", intent.Id)
	log.Info("received scraping intent")

	limit <- true
	client := scraper.NewScraperTasksServiceClient(runner.conn)

	stub := empty.Empty{}
	task, err := client.CreateTask(context, &stub)
	if err != nil {
		<-limit
		log.Errorf("failed to create scraping task: %v", err)
		return err
	}

	go func() {
		runScraper(scrapeContext{
			intent: intent,
			task:   task,
			client: client,
			log:    log.WithField("task", task.Id),
		})
		<-limit
	}()

	return nil
}

type scrapeContext struct {
	intent *scraper.ScrapeIntent
	task   *scraper.Task
	client scraper.ScraperTasksServiceClient
	log    *logrus.Entry
}

func runScraper(ctx scrapeContext) {
	switch ctx.intent.Source {
	case scraper.ScrapeIntent_ANIDB:
		startAniDBScraping(ctx)
	}
}

func startAniDBScraping(ctx scrapeContext) {
	proxies := runner.proxyFetcher.Fetch()

	tr := grpcTransport{client: ctx.client}
	reporter := NewTaskReporter(ctx.task, tr)
	spider := anidb.NewSpider(ctx.task, reporter)

	if len(proxies) == 0 {
		reporter.Finish()
		return
	}

	spider.SetProxies(proxies)
	spider.Run()
}

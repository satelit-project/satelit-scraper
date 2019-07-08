package anidb

import (
	"bytes"
	"fmt"
	"runtime"

	"satelit-project/satelit-scraper/proto/scraper"
	"satelit-project/satelit-scraper/spider"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/proxy"
	"github.com/gocolly/colly/queue"
	"github.com/sirupsen/logrus"
)

type Spider struct {
	task *scraper.Task
	reporter *spider.TaskReporter
	proxies []string
	urlMap map[string]int32
	log *logrus.Entry
}

func (s *Spider) Run() {
	coll := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
		colly.Debugger(CollyLogger{s.log}),
	)

	s.setupProxy(coll)
	s.setupCallbacks(coll)

	collq, err := s.makeQueue()
	if err != nil {
		s.log.Errorf("failed to setup queue: %v", err)
		return
	}

	if err := collq.Run(coll); err != nil {
		s.log.Errorf("scraping failed: %v", err)
	}
}

func (s *Spider) setupProxy(coll *colly.Collector) {
	prx, err := proxy.RoundRobinProxySwitcher(s.proxies...)
	if err != nil {
		s.log.Errorf("failed to set up proxy: %v", err)
		return
	}

	coll.SetProxyFunc(prx)
}

func (s *Spider) setupCallbacks(coll *colly.Collector) {
	coll.OnResponse(func(r *colly.Response) {
		parser, err := NewParser(r.Request.URL, bytes.NewReader(r.Body))
		if err != nil {
			s.log.Errorf("failed to create parser: %v", err)
			return
		}

		anime, err := parser.Anime()
		if err != nil {
			s.log.Errorf("failed to parse anime: %v", err)
			return
		}

		scheduleID, ok := s.urlMap[r.Request.URL.String()]
		if ok != true {
			s.log.Errorf("schedule id not found")
		}

		s.reporter.Report(anime, scheduleID)
	})
}

func (s *Spider) makeQueue() (*queue.Queue, error) {
	collq, err := queue.New(
		runtime.NumCPU(),
		&queue.InMemoryQueueStorage{MaxSize: len(s.task.ScheduleIds)},
	)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(s.task.AnimeIds); i++ {
		url := urlForID(s.task.AnimeIds[i])
		err := collq.AddURL(url)
		s.urlMap[url] = s.task.ScheduleIds[i]
		if err != nil {
			// will not happen
			s.log.Errorf("failed to queue url: %v", err)
		}
	}

	return collq, nil
}

func urlForID(id int32) string {
	return fmt.Sprintf("https://anidb.net/perl-bin/animedb.pl?show=anime&aid=%d", id)
}

type CollyLogger struct {
	log *logrus.Entry
}

func (l CollyLogger) Init() error {
	return nil
}

func (l CollyLogger) Event(e *debug.Event) {
	l.log.Debugf("%d [%6d - %s] %q (%s)\n", e.CollectorID, e.RequestID, e.Type, e.Values)
}

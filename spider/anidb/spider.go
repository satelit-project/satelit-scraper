package anidb

import (
	"bytes"
	"fmt"
	"time"

	"satelit-project/satelit-scraper/proto/scraper"
	"satelit-project/satelit-scraper/spider"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"github.com/sirupsen/logrus"
)

type Spider struct {
	task     *scraper.Task
	reporter *spider.TaskReporter
	proxies  []string
	urlMap   map[string]int32
	timeout  time.Duration
	delay    time.Duration
	log      *logrus.Entry
}

func NewSpider(task *scraper.Task, reporter *spider.TaskReporter) *Spider {
	log := logrus.WithField("spider_task", task.Id)

	return &Spider{
		task:     task,
		reporter: reporter,
		proxies:  make([]string, 0),
		urlMap:   make(map[string]int32),
		timeout:  45 * time.Second,
		delay:    3 * time.Second,
		log:      log,
	}
}

func (s *Spider) SetProxies(proxies []string) {
	s.proxies = proxies
}

func (s *Spider) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *Spider) SetDelay(delay time.Duration) {
	s.delay = delay
}

func (s *Spider) Run() {
	coll := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
		colly.Debugger(CollyLogger{s.log}),
	)

	s.setupProxy(coll)
	s.setupCallbacks(coll)
	coll.SetRequestTimeout(s.timeout)
	_ = coll.Limit(&colly.LimitRule{DomainGlob: "*", Delay: s.delay})

	animeURLs := s.makeURLs()
	for _, animeURL := range animeURLs {
		err := coll.Visit(animeURL)
		if err != nil {
			s.log.Errorf("scraping failed: %v", err)
		}
	}

	coll.Wait()
	s.reporter.Finish()
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
	extensions.RandomUserAgent(coll)

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
		if !ok {
			s.log.Errorf("schedule id not found")
		}

		s.reporter.Report(anime, scheduleID)
	})

	coll.OnError(func(r *colly.Response, e error) {
		s.log.Errorf("request failed: %v", e)
	})
}

func (s *Spider) makeURLs() []string {
	var urls []string
	for i := 0; i < len(s.task.AnimeIds); i++ {
		url := urlForID(s.task.AnimeIds[i])
		urls = append(urls, url)
		s.urlMap[url] = s.task.ScheduleIds[i]
	}

	return urls
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
	l.log.Debugf("%d [%6d - %s] %q\n", e.CollectorID, e.RequestID, e.Type, e.Values)
}

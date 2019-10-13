package anidb

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	cproxy "github.com/gocolly/colly/proxy"
	"go.uber.org/zap"

	"satelit-project/satelit-scraper/logging"
	"satelit-project/satelit-scraper/proto/scraping"
	"satelit-project/satelit-scraper/proxy"
	"satelit-project/satelit-scraper/spider"
)

type Spider struct {
	task     *scraping.Task
	reporter *spider.TaskReporter
	proxies  []proxy.Proxy
	jobMap   map[string]int
	timeout  time.Duration
	delay    time.Duration
	log      *zap.SugaredLogger
}

func NewSpider(task *scraping.Task, reporter *spider.TaskReporter) *Spider {
	log := logging.DefaultLogger().With("spider_task", task.Id)

	return &Spider{
		task:     task,
		reporter: reporter,
		proxies:  make([]proxy.Proxy, 0),
		jobMap:   make(map[string]int),
		timeout:  20 * time.Second,
		delay:    3 * time.Second,
		log:      log,
	}
}

func (s *Spider) SetProxies(proxies []proxy.Proxy) {
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
	coll.DisableCookies()
	extensions.RandomUserAgent(coll)
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
	proxies := make([]string, 0, len(s.proxies))
	for _, p := range s.proxies {
		proxies = append(proxies, p.String())
	}

	prx, err := cproxy.RoundRobinProxySwitcher(proxies...)
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

		jobIndex, ok := s.jobMap[r.Request.URL.String()]
		if !ok {
			s.log.Errorf("job index not found")
			return
		}

		job := s.task.Jobs[jobIndex]
		s.reporter.Report(job, anime)
	})

	coll.OnRequest(func(r *colly.Request) {
		if r.ProxyURL != "" {
			s.log.Infof("using proxy: %s", r.ProxyURL)
		}
	})

	coll.OnError(func(r *colly.Response, e error) {
		s.log.Errorf("request failed: %v", e)
	})
}

func (s *Spider) makeURLs() []string {
	var urls []string
	for i := 0; i < len(s.task.Jobs); i++ {
		url := urlForID(s.task.Jobs[i].AnimeId)
		urls = append(urls, url)
		s.jobMap[url] = i
	}

	return urls
}

func urlForID(id int32) string {
	return fmt.Sprintf("https://anidb.net/perl-bin/animedb.pl?show=anime&aid=%d", id)
}

type CollyLogger struct {
	log *zap.SugaredLogger
}

func (l CollyLogger) Init() error {
	return nil
}

func (l CollyLogger) Event(e *debug.Event) {
	l.log.Debugf("%d [%6d - %s] %q\n", e.CollectorID, e.RequestID, e.Type, e.Values)
}

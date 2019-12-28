package anidb

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	cproxy "github.com/gocolly/colly/proxy"

	"shitty.moe/satelit-project/satelit-scraper/config"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/parser/anidb"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
	"shitty.moe/satelit-project/satelit-scraper/proxy"
	"shitty.moe/satelit-project/satelit-scraper/spider"
)

// Performs AniDB scraping task.
type Spider struct {
	task     *scraping.Task
	reporter *spider.TaskReporter
	config   *config.AniDB
	proxies  []proxy.Proxy
	jobMap   map[string]int
	log      *logging.Logger
}

// Creates and returns new Spider instance.
func NewSpider(reporter *spider.TaskReporter, config *config.AniDB, log *logging.Logger) Spider {
	log = log.With("anidb_task", reporter.Task.Id)

	return Spider{
		task:     reporter.Task,
		reporter: reporter,
		config:   config,
		proxies:  nil,
		jobMap:   make(map[string]int),
		log:      log,
	}
}

// Sets list of proxy servers to use when making HTTP requests.
func (s *Spider) SetProxies(proxies []proxy.Proxy) {
	s.proxies = proxies
}

// Starts scraping process.
func (s *Spider) Run() {
	coll := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
		colly.Debugger(logging.CollyLogger{Log: s.log}),
	)

	s.setupProxy(coll)
	s.setupCallbacks(coll)
	coll.SetRequestTimeout(time.Duration(s.config.Timeout) * time.Second)
	coll.DisableCookies()
	extensions.RandomUserAgent(coll)
	_ = coll.Limit(&colly.LimitRule{DomainGlob: "*", Delay: time.Duration(s.config.Delay) * time.Second})

	animeURLs := s.makeURLs()
	for _, animeURL := range animeURLs {
		err := coll.Visit(animeURL)
		if err != nil {
			s.log.Errorf("scraping failed: %v", err)
		}
	}

	coll.Wait()
	if err := s.reporter.Finish(); err != nil {
		s.log.Errorf("failed to report scraping finished: %v", err)
	}
}

// Setups proxy for the scraper.
func (s *Spider) setupProxy(coll *colly.Collector) {
	if len(s.proxies) == 0 {
		return
	}

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

// Setups Colly callbacks for the spider.
func (s *Spider) setupCallbacks(coll *colly.Collector) {
	coll.OnResponse(func(r *colly.Response) {
		parser, err := anidb.NewParser(r.Request.URL, bytes.NewReader(r.Body), s.log)
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
		if err = s.reporter.Report(job, anime); err != nil {
			s.log.Errorf("failed to report job: %v", err)
		}
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

// Makes list of AniDB URLs to visit for parsing. The method will also fill scraper's jobMap property.
func (s *Spider) makeURLs() []string {
	var urls []string
	for i := 0; i < len(s.task.Jobs); i++ {
		id := s.task.Jobs[i].AnimeId
		url := fmt.Sprintf(s.config.URLTemplate, id)
		urls = append(urls, url)
		s.jobMap[url] = i
	}

	return urls
}

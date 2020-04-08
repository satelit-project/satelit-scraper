package main

import (
	"context"
	"flag"
	"time"

	"shitty.moe/satelit-project/satelit-scraper/config"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/proto/scraping"
	"shitty.moe/satelit-project/satelit-scraper/server"

	uuid "shitty.moe/satelit-project/satelit-scraper/proto/common"
	"shitty.moe/satelit-project/satelit-scraper/proto/data"
)

var port = flag.Int("port", 10602, "port to use for grpc service")
var animeIDs = AnimeIDsFlag("id", nil, "anime id to scrape")

func main() {
	flag.Parse()

	log, err := logging.NewLogger(nil)
	if err != nil {
		panic(err)
	}

	cfg, err := config.Default()
	if err != nil {
		panic(err)
	}

	srvc := NewTaskService(*port, log)
	go func() {
		var jobs []int32
		for _, j := range *animeIDs {
			jobs = append(jobs, j)
		}

		srvc.SetJobs(jobs)
		if err := srvc.Run(); err != nil {
			panic(err)
		}
	}()

	if err = WaitPort(*port, 10*time.Second); err != nil {
		panic(err)
	}

	runner := server.NewRunner(cfg.Scraping, cfg.AniDB, noCache{}, log)
	_, err = runner.Run(context.Background(), &scraping.ScrapeIntent{
		Id:     &uuid.UUID{},
		Source: data.Source_ANIDB,
	})

	if err != nil {
		panic(err)
	}
}

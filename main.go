package main

import (
	"os"
	"os/signal"
	"syscall"

	"shitty.moe/satelit-project/satelit-scraper/config"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/server"
	"shitty.moe/satelit-project/satelit-scraper/spider"
)

func main() {
	cfg := makeConfig()
	log := makeLogger(cfg)
	defer func() {
		_ = log.Sync()
	}()

	cache, err := makeCache(cfg.Storage, log)
	if err != nil {
		log.Errorf("failed to create cache: %v", err)
		return
	}

	srvc := server.New(cfg, cache, log)
	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-done

		srvc.Shutdown()
	}()

	log.Infof("starting service")
	if err := srvc.Run(); err != nil {
		log.Errorf("error while running service: %v", err)
	}

	log.Infof("service stopped")
}

func makeLogger(cfg config.Config) *logging.Logger {
	log, err := logging.NewLogger(cfg.Logging)
	if err != nil {
		panic(err)
	}

	return log
}

func makeConfig() config.Config {
	cfg, err := config.Default()
	if err != nil {
		panic(err)
	}

	return cfg
}

func makeCache(cfg *config.Storage, log *logging.Logger) (spider.Cache, error) {
	return spider.NewS3Cache(cfg, log)
}

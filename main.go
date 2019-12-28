package main

import (
	"os"
	"os/signal"
	"syscall"

	"shitty.moe/satelit-project/satelit-scraper/config"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/server"
)

func main() {
	cfg := makeConfig()
	log := makeLogger()
	defer func() {
		_ = log.Sync()
	}()

	srvc := server.New(cfg, log)

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

func makeLogger() *logging.Logger {
	log, err := logging.NewLogger()
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

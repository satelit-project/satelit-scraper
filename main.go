package main

import "os"
import "os/signal"
import "syscall"

import "shitty.moe/satelit-project/satelit-scraper/server"

import "shitty.moe/satelit-project/satelit-scraper/config"

import "shitty.moe/satelit-project/satelit-scraper/logging"

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

	if err := srvc.Run(); err != nil {
		log.Errorf("error while running service: %v", err)
	}
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

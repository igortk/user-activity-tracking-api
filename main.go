package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"user-activity-tracking-api/config"
	"user-activity-tracking-api/cron"
	"user-activity-tracking-api/rest"
	"user-activity-tracking-api/service/database"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		log.Fatalf("Failed load configuration: %v", err)
	}

	database.ConnectDB(&cfg.DataBaseConfig)
	startServers(cfg)
}

func startServers(cfg *config.Config) {
	var wg sync.WaitGroup

	stopCh := make(chan struct{})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg.Add(2)

	go rest.Run(&wg, &cfg.HttpConfig, stopCh)
	go cron.Run(&wg, &cfg.CronConfig, stopCh)

	select {
	case s := <-sigCh:
		log.Infof("Received OS signal: %v", s)
	}

	close(stopCh)

	wg.Wait()
}

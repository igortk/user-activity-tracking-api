package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"user-activity-tracking-api/config"
	"user-activity-tracking-api/cron"
	"user-activity-tracking-api/rest"
	"user-activity-tracking-api/service/database"
	"user-activity-tracking-api/service/database/aggregators"
)

func main() {
	cfg, err := config.GetConfig()

	initSetOutputLogs()

	if err != nil {
		log.Fatalf("Failed load configuration: %v", err)
	}
	ctx := context.Background()
	db := database.NewClient(&cfg.DataBaseConfig)
	startServers(ctx, cfg, db)
}

func initSetOutputLogs() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
}

func startServers(ctx context.Context, cfg *config.Config, dbCl *database.Client) {
	var wg sync.WaitGroup

	stopCh := make(chan struct{})

	//var agr aggregators.UserEventsAggregator = aggregators.NewSQLUserEventsAggregator(dbCl.GetDb())
	var agr aggregators.UserEventsAggregator
	agr = aggregators.NewSQLUserEventsAggregator(dbCl.GetDb())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg.Add(2)

	runRestServer(&wg, cfg, dbCl, stopCh)
	runCronServer(ctx, &wg, agr, dbCl, cfg, stopCh)

	select {
	case s := <-sigCh:
		log.Infof("Received OS signal: %v", s)

		log.Info("Closing db connection...")
		dbCl.Close()
	}

	close(stopCh)

	wg.Wait()
}

func runRestServer(wg *sync.WaitGroup, cfg *config.Config, dbCl *database.Client, stopCh <-chan struct{}) {
	srv := rest.NewServer(cfg, dbCl)

	go srv.Run(wg, stopCh)
}

func runCronServer(ctx context.Context, wg *sync.WaitGroup, agr aggregators.UserEventsAggregator,
	dbCl *database.Client, cfg *config.Config, stopCh <-chan struct{}) {
	srv := cron.NewServer(ctx, cfg, agr, dbCl) //TODO agr

	go srv.Run(wg, stopCh)
}

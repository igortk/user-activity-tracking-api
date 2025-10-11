package cron

import (
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	"user-activity-tracking-api/config"
	"user-activity-tracking-api/cron/jobs"
)

func Run(wg *sync.WaitGroup, cronCfg *config.CronConfig, stopCh <-chan struct{}) {
	defer wg.Done()

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Cron(cronCfg.TabCountUsersEventTask).Do(jobs.CalculateUserEventsAndSaveDb)
	if err != nil {
		log.Errorf("Failed to schedule cron testTask: %v", err)
		return
	}

	s.StartAsync()

	<-stopCh
	log.Info("Stop signal received. Stopping cron scheduler...")

	s.Stop()
	log.Info("Cron server stopped")
}

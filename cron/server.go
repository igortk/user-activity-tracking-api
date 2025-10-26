package cron

import (
	"context"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"time"
	"user-activity-tracking-api/config"
	"user-activity-tracking-api/cron/jobs"
	"user-activity-tracking-api/service/database"
	"user-activity-tracking-api/service/database/aggregators"
	"user-activity-tracking-api/service/database/repositories"
)

type Server struct {
	cronCfg             *config.CronConfig
	aggregator          aggregators.UserEventsAggregator
	userEventCountsRepo *repositories.UserEventCountsRepository

	cronTab []Tab
}

func NewServer(ctx context.Context, cfg *config.Config, agr aggregators.UserEventsAggregator, dbCl *database.Client) *Server {
	numCronTab := reflect.TypeOf(cfg.CronConfig.Tab).NumField()

	srv := &Server{
		cronCfg:             &cfg.CronConfig,
		aggregator:          agr,
		userEventCountsRepo: repositories.NewUserEventCountsRepository(dbCl.GetDb()),

		cronTab: make([]Tab, 0, numCronTab),
	}

	srv.initCronTabs(ctx)

	return srv
}

func (s *Server) Run(wg *sync.WaitGroup, stopCh <-chan struct{}) {
	defer wg.Done()

	srv := gocron.NewScheduler(time.UTC)

	s.initCron(srv)

	srv.StartAsync()

	<-stopCh
	log.Info("Stop signal received. Stopping cron scheduler...")

	srv.Stop()
	log.Info("Cron server stopped")
}

func (s *Server) initCronTabs(ctx context.Context) {
	f := jobs.NewCalculateUserEventsAndSaveDb(ctx, s.aggregator, s.userEventCountsRepo)

	s.cronTab = append(s.cronTab, Tab{
		Schedule: s.cronCfg.Tab.TabCountUsersEventTask,
		Job:      f,
	})
}

func (s *Server) initCron(srv *gocron.Scheduler) {
	for _, tab := range s.cronTab {
		if _, err := srv.Cron(tab.Schedule).Do(tab.Job); err != nil {
			log.Errorf("Failed to schedule cron testTask: %v", err)
		}
	}
}

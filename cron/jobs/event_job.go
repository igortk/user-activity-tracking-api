package jobs

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
	"user-activity-tracking-api/service/database/aggregators"
	"user-activity-tracking-api/service/database/repositories"
)

type CalculateUserEventsAndSaveDb struct {
	ctx context.Context

	userEventsRepo         aggregators.UserEventsAggregator
	userEventCountsRepo    *repositories.UserEventCountsRepository
	startCalculationPeriod time.Time
}

func NewCalculateUserEventsAndSaveDb(ctx context.Context, agr aggregators.UserEventsAggregator,
	repo *repositories.UserEventCountsRepository) *CalculateUserEventsAndSaveDb {

	return &CalculateUserEventsAndSaveDb{
		ctx: ctx,

		userEventCountsRepo:    repo,
		userEventsRepo:         agr,
		startCalculationPeriod: time.Now().UTC(),
	}
}

func (c *CalculateUserEventsAndSaveDb) Run() {
	log.Infof("Start calculate user events, start period: %s", c.startCalculationPeriod)

	currentTime := time.Now().UTC()
	evCounts, err := c.userEventsRepo.AggregateUserEvents(c.ctx, c.startCalculationPeriod, currentTime)

	if err != nil {
		log.Errorf("Failed to count user event: %v", err)
	} else {
		err = c.userEventCountsRepo.SaveUserEvents(evCounts)
		if err != nil {
			log.Errorf("Failed to save counted user event: %v", err)
		}
	}

	c.startCalculationPeriod = currentTime
}

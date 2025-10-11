package jobs

import (
	log "github.com/sirupsen/logrus"
	"time"
	"user-activity-tracking-api/service/database/repositories"
)

var startCalculationPeriod time.Time

func init() {
	startCalculationPeriod = time.Now()
}

func CalculateUserEventsAndSaveDb() {
	log.Infof("Start calculate user events, start period: %s", startCalculationPeriod)

	currentTime := time.Now()
	err := repositories.CalculateAndSaveUserEvents(startCalculationPeriod, currentTime)
	if err != nil {
		log.Error("Failed to count user event")
	}

	startCalculationPeriod = currentTime
}

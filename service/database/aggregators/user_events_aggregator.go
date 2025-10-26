package aggregators

import (
	"context"
	"time"
	"user-activity-tracking-api/models"
)

type UserEventsAggregator interface {
	AggregateUserEvents(ctx context.Context, start, end time.Time) ([]models.UserEventCount, error)
}

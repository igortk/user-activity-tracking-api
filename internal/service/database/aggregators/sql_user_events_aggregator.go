package aggregators

import (
	"context"
	"gorm.io/gorm"
	"time"
	"user-activity-tracking-api/internal/models"
)

type SQLUserEventsAggregator struct {
	db *gorm.DB
}

func NewSQLUserEventsAggregator(db *gorm.DB) *SQLUserEventsAggregator {
	return &SQLUserEventsAggregator{
		db: db,
	}
}

func (a *SQLUserEventsAggregator) AggregateUserEvents(ctx context.Context, start, end time.Time) ([]models.UserEventCount, error) {
	var results []models.UserEventCount
	query := `
        SELECT 
            user_id,
            $1::timestamptz AS period_start,
            COUNT(*) AS event_count,
            $2::timestamptz AS period_end
        FROM events
        WHERE created_at BETWEEN $1 AND $2
        GROUP BY user_id
    `
	err := a.db.WithContext(ctx).Raw(query, start, end).Scan(&results).Error
	return results, err
}

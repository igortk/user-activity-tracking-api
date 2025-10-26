package repositories

import (
	"context"
	"gorm.io/gorm"
	"time"
	"user-activity-tracking-api/models"
)

type EventsRepository struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) *EventsRepository {
	return &EventsRepository{
		db: db,
	}
}

func (r *EventsRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *EventsRepository) GetEventsByUserIdAndDateRange(ctx context.Context, userID, limit, offset int64, from, to time.Time) ([]models.Event, error) {
	var events []models.Event

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND event_action_timestamp BETWEEN ? AND ?", userID, from, to).
		Order("event_action_timestamp ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&events)

	return events, result.Error
}

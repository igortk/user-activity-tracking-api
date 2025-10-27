package repositories

import (
	"context"
	"gorm.io/gorm"
	"time"
	"user-activity-tracking-api/internal/models"
)

type EventsRepository interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEventsByUserIdAndDateRange(ctx context.Context, userID, limit, offset int64, from, to time.Time) ([]models.Event, error)
}

type GormEventsRepository struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) *GormEventsRepository {
	return &GormEventsRepository{
		db: db,
	}
}

func (r *GormEventsRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *GormEventsRepository) GetEventsByUserIdAndDateRange(ctx context.Context, userID, limit, offset int64, from, to time.Time) ([]models.Event, error) {
	var events []models.Event

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND event_action_timestamp BETWEEN ? AND ?", userID, from, to).
		Order("event_action_timestamp ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&events)

	return events, result.Error
}

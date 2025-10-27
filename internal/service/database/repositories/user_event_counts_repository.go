package repositories

import (
	"gorm.io/gorm"
	"user-activity-tracking-api/internal/models"
)

type UserEventCountsRepository struct {
	db *gorm.DB
}

func NewUserEventCountsRepository(db *gorm.DB) *UserEventCountsRepository {
	return &UserEventCountsRepository{
		db: db,
	}
}

func (s *UserEventCountsRepository) SaveUserEvents(eventsCount []models.UserEventCount) error {
	if len(eventsCount) < 1 {
		return nil
	}
	result := s.db.Create(&eventsCount)
	return result.Error
}

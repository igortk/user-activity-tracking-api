package repositories

import (
	"gorm.io/gorm"
	"user-activity-tracking-api/models"
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
	result := s.db.Create(&eventsCount)
	return result.Error
}

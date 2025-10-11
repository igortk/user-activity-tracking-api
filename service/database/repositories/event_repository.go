package repositories

import (
	"time"
	"user-activity-tracking-api/models"
	"user-activity-tracking-api/service/database"
)

func GetEventsByUserIdAndDateRange(userID int, from, to time.Time) ([]models.Event, error) {
	var events []models.Event

	result := database.Session.
		Where("user_id = ? AND event_action_timestamp BETWEEN ? AND ?", userID, from, to).
		Order("event_action_timestamp ASC").
		Find(&events)

	return events, result.Error
}

package models

import "time"

type UserEventCount struct {
	UserID      uint      `json:"user_id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	EventCount  int64     `json:"event_count"`
}

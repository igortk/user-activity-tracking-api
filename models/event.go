package models

import (
	"encoding/json"
	"time"
)

type Event struct {
	UserID               int64           `json:"user_id" validate:"required,gt=0"`
	EventActionTimestamp time.Time       `json:"event_action_timestamp" validate:"required"`
	Action               string          `json:"action" validate:"required,oneof=created updated deleted viewed"`
	Metadata             json.RawMessage `json:"metadata" validate:"required,json"`
}

package models

import "time"

type GetActivityEventByUserIdDateRangeRequest struct {
	UserID                   int       `form:"user_id" validate:"required,gt=0"`
	FromEventActionTimestamp time.Time `form:"from" validate:"required"`
	ToEventActionTimestamp   time.Time `form:"to" validate:"required"`
}

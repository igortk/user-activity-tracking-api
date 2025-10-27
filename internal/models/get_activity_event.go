package models

import "time"

type GetActivityEventByUserIdDateRangeRequest struct {
	UserID                   int64     `form:"user_id" validate:"required,gt=0"`
	FromEventActionTimestamp time.Time `form:"from" validate:"required"`
	ToEventActionTimestamp   time.Time `form:"to" validate:"required"`
	Offset                   int64     `form:"offset" validate:"required"`
	Limit                    int64     `form:"limit" validate:"required"`
}

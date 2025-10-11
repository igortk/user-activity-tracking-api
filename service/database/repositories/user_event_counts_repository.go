package repositories

import (
	"time"
	"user-activity-tracking-api/service/database"
)

func CalculateAndSaveUserEvents(startPeriod, endPeriod time.Time) error {
	sql := `
		INSERT INTO user_event_counts (user_id, period_start, event_count, period_end)
		SELECT 
			user_id,
			$1::timestamptz AS period_start,
			COUNT(*) AS event_count,
		    $2::timestamptz AS period_end
		FROM events
		WHERE created_at BETWEEN $1::timestamptz AND $2::timestamptz
		GROUP BY user_id
	`

	result := database.Session.Exec(sql, startPeriod, endPeriod)
	return result.Error
}

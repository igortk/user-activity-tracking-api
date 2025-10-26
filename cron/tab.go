package cron

import (
	"github.com/robfig/cron/v3"
)

type Tab struct {
	Schedule string
	Job      cron.Job
}

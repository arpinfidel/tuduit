package entity

import (
	"time"
)

type Schedule struct {
	StdFields

	StartDate time.Time `json:"start_date" db:"start_date"`
	Schedule  string    `json:"schedule" db:"schedule"` // cron expression
	Duration  int64     `json:"duration" db:"duration"` // duration before deadline in seconds
}

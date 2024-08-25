package entity

import "time"

type CheckIn struct {
	StdFields

	UserID      int64     `db:"user_id"       json:"user_id"`
	CheckInTime time.Time `db:"check_in_time" json:"check_in_time"`
	LastSent    time.Time `db:"last_sent"     json:"last_sent"`
}

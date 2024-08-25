package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
	str2duration "github.com/xhit/go-str2duration/v2"
)

type Schedule struct {
	StdFields

	UserID int64 `db:"user_id" json:"user_id"`

	StartDate    time.Time  `json:"start_date" db:"start_date"`
	EndDate      *time.Time `json:"end_date" db:"end_date"`
	Schedule     string     `json:"schedule" db:"schedule"` // cron expression
	NextSchedule time.Time  `json:"next_schedule" db:"next_schedule"`
	Duration     int64      `json:"duration" db:"duration"` // duration before deadline in seconds
}

type ScheduleExprType string

const (
	ScheduleExprTypeCron ScheduleExprType = "cron"
	ScheduleExprTypeFreq ScheduleExprType = "freq"
)

type ScheduleExpr struct {
	Type ScheduleExprType
	Cron string
	Freq int64 // in seconds
}

func (s *Schedule) ParseSchedule() (ScheduleExpr, error) {
	if s.Schedule == "" {
		return ScheduleExpr{}, nil
	}

	errs := []string{}

	_, err := cronexpr.Parse(s.Schedule)
	if err == nil {
		return ScheduleExpr{
			Type: ScheduleExprTypeCron,
			Cron: s.Schedule,
		}, nil
	}
	errs = append(errs, fmt.Errorf("invalid cron format: %s", s.Schedule).Error())

	freq, err := strconv.ParseInt(s.Schedule, 10, 64)
	if err == nil {
		return ScheduleExpr{
			Type: ScheduleExprTypeFreq,
			Freq: freq,
		}, nil
	}
	errs = append(errs, fmt.Errorf("invalid frequency format: %s", s.Schedule).Error())

	secs, err := str2duration.ParseDuration(s.Schedule)
	if err == nil {
		return ScheduleExpr{
			Type: ScheduleExprTypeFreq,
			Freq: int64(secs.Seconds()),
		}, nil
	}
	errs = append(errs, fmt.Errorf("invalid value for 1w1d1h1m1s format: %w", err).Error())

	return ScheduleExpr{}, fmt.Errorf("invalid schedule: %s\nerrors: %s", s.Schedule, strings.Join(errs, "\n"))
}

func (s *Schedule) MustParseSchedule() *ScheduleExpr {
	expr, err := s.ParseSchedule()
	if err != nil {
		panic(err)
	}
	return &expr
}

func (s *ScheduleExpr) Next(now time.Time) (time.Time, error) {
	switch s.Type {
	case ScheduleExprTypeCron:
		cron, err := cronexpr.Parse(s.Cron)
		if err != nil {
			return time.Time{}, err
		}
		return cron.Next(now), nil
	case ScheduleExprTypeFreq:
		return now.Add(time.Duration(s.Freq) * time.Second), nil
	}

	return time.Time{}, fmt.Errorf("invalid schedule type: %s", s.Type)
}

func (s *ScheduleExpr) MustNext(now time.Time) time.Time {
	next, err := s.Next(now)
	if err != nil {
		panic(err)
	}
	return next
}

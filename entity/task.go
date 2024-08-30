package entity

import (
	"fmt"
	"time"
)

type Task struct {
	StdFields

	UserIDs        PQArr[int64] `db:"user_ids"         json:"user_ids"`
	TaskScheduleID int64        `db:"task_schedule_id" json:"task_schedule_id"`

	Name        string     `db:"name"         json:"name"`
	Description string     `db:"description"  json:"description"`
	Priority    int        `db:"priority"     json:"priority"`
	Status      string     `db:"status"       json:"status"`
	StartDate   *time.Time `db:"start_date"   json:"start_date,omitempty"`
	EndDate     *time.Time `db:"end_date"     json:"end_date,omitempty"`
	StartedAt   *time.Time `db:"started_at"   json:"started_at,omitempty"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	ArchivedAt  *time.Time `db:"archived_at"  json:"archived_at,omitempty"`
	IsTemplate  bool       `db:"is_template"  json:"is_template"`
}

func (t Task) Overview() TaskOverview {
	to := TaskOverview{
		ID:          t.ID,
		Name:        t.Name,
		Priority:    t.Priority,
		Description: t.Description,
		Status:      t.Status,
	}

	if t.CompletedAt != nil {
		to.Completion = "completed"
	} else if t.StartedAt != nil {
		to.Completion = "in progress"
	} else {
		to.Completion = "not started"
	}

	if t.StartDate != nil {
		to.StartDate = t.StartDate.Format("2006-01-02 15:04:05")
		elapsed := Duration(time.Since(*t.StartDate).Round(time.Second))
		if elapsed < 0 {
			to.StartDate += fmt.Sprintf(" (in %s)", -elapsed)
		} else {
			to.StartDate += fmt.Sprintf(" (%s ago)", elapsed)
		}
	}
	if t.EndDate != nil {
		to.EndDate = t.EndDate.Format("2006-01-02 15:04:05")
		elapsed := Duration(time.Since(*t.EndDate).Round(time.Second))
		if elapsed < 0 {
			to.EndDate += fmt.Sprintf(" (in %s)", -elapsed)
		} else {
			to.EndDate += fmt.Sprintf(" (%s ago)", elapsed)
		}
	}

	return to
}

type TaskOverview struct {
	ID int64 `db:"id" json:"id"`

	Name        string `db:"name"        json:"name"                 yaml:"name"`
	Priority    int    `db:"priority"    json:"priority"             yaml:"priority"`
	Description string `db:"description" json:"description"          yaml:"description,omitempty"`
	Status      string `db:"status"      json:"status,omitempty"     yaml:"status,omitempty"`
	StartDate   string `db:"start_date"  json:"start_date,omitempty" yaml:"start_date,omitempty"`
	EndDate     string `db:"end_date"    json:"end_date,omitempty"   yaml:"end_date,omitempty"`

	Completion string `json:"completion"`
}

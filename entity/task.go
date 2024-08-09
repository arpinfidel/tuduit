package entity

import "time"

type Task struct {
	StdFields

	UserID int `db:"user_id" json:"user_id"`

	Name        string     `db:"name"         json:"name"`
	Description string     `db:"description"  json:"description"`
	Priority    int        `db:"priority"     json:"priority"`
	Status      string     `db:"status"       json:"status"`
	StartedAt   *time.Time `db:"started_at"   json:"started_at,omitempty"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	ArchivedAt  *time.Time `db:"archived_at"  json:"archived_at,omitempty"`
}

func (t Task) Overview() TaskOverview {
	completion := ""
	if t.CompletedAt != nil {
		completion = "completed"
	} else if t.StartedAt != nil {
		completion = "in progress"
	} else {
		completion = "not started"
	}
	return TaskOverview{
		ID:          t.ID,
		Name:        t.Name,
		Priority:    t.Priority,
		Description: t.Description,
		Status:      t.Status,

		Completion: completion,
	}
}

type TaskOverview struct {
	ID int `db:"id" json:"id"`

	Name        string `db:"name"        json:"name"`
	Priority    int    `db:"priority"    json:"priority"`
	Description string `db:"description" json:"description"`
	Status      string `db:"status"      json:"status"`

	Completion string `json:"completion"`
}

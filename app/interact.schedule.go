package app

import (
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
)

type CreateScheduleParams struct {
	StartDate time.Time `rose:"start_date,s,required="`
	Schedule  string    `rose:"schedule,sc,required="` // cron expression

	Name     string `rose:"name,n,required="`
	Priority int    `rose:"priority,p,default=2"`

	Description string     `rose:"description,d"`
	EndDate     *time.Time `rose:"end_date,e"`
	Assignee    string     `rose:"assignee,a"`
}

func (h *App) CreateSchedule(ctx *ctxx.Context, p CreateScheduleParams) (next string, err error) {

	userID := ctx.UserID
	if p.Assignee != "" {
		user, _, err := h.d.UserUC.Get(ctx, nil, db.Params{
			Where: []db.Where{
				{
					Field: "username",
					Op:    db.EqOp,
					Value: p.Assignee,
				},
			},
		})
		if err != nil {
			return "", err
		}
		if len(user) == 0 {
			return "", fmt.Errorf("user not found")
		}

		userID = user[0].ID
	}

	schedule := entity.Schedule{
		UserID:    userID,
		StartDate: p.StartDate,
		EndDate:   p.EndDate,
		Schedule:  p.Schedule,
	}

	expr, err := schedule.ParseSchedule()
	if err != nil {
		return "", err
	}
	switch expr.Type {
	case entity.ScheduleExprTypeCron:
		schedule.NextSchedule = expr.MustNext(p.StartDate)
	case entity.ScheduleExprTypeFreq:
		schedule.NextSchedule = p.StartDate
	}

	s, err := h.d.ScheduleUC.Create(ctx, nil, []entity.Schedule{schedule})
	if err != nil {
		return "", err
	}
	sched := s[0]

	template := entity.Task{
		UserID:         userID,
		TaskScheduleID: sched.ID,
		Name:           p.Name,
		Priority:       p.Priority,
		Description:    p.Description,
		IsTemplate:     true,
	}

	task := template
	task.IsTemplate = false
	task.StartDate = &schedule.NextSchedule

	_, err = h.d.TaskUC.Create(ctx, nil, []entity.Task{
		template,
		task,
	})
	if err != nil {
		return "", err
	}

	return schedule.NextSchedule.Format("2006-01-02 15:04:05"), nil
}

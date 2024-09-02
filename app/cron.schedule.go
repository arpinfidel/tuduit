package app

import (
	"context"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
)

var _ = registerTask("CreateScheduledTasks", "* * * * *", func() func() error { return a.CreateScheduledTasks })

func (a *App) CreateScheduledTasks() error {
	ctx := context.Background()

	now := time.Now().UTC()

	scheds, _, err := a.d.ScheduleUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "next_schedule",
				Op:    db.LtOrEqOp,
				Value: now,
			},
		},
		Pagination: &db.Pagination{
			Limit:  1 << 30,
			Offset: 0,
		},
	})
	if err != nil {
		return err
	}
	if len(scheds) == 0 {
		return nil
	}
	schedIDs := []int64{}
	for _, s := range scheds {
		schedIDs = append(schedIDs, s.ID)
	}

	templates, _, err := a.d.TaskUC.GetTemplate(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "task_schedule_id",
				Op:    db.InOp,
				Value: schedIDs,
			},
		},
	})
	if err != nil {
		return err
	}

	mapNextScheds := map[int64]entity.Schedule{}
	for i := range scheds {
		s := &scheds[i]
		s.NextSchedule = s.MustParseSchedule().MustNext(s.NextSchedule)
		mapNextScheds[s.ID] = *s
		_, err = a.d.ScheduleUC.Update(ctxx.New(ctx, entity.User{}), nil, *s)
		if err != nil {
			return err
		}
	}

	tasks := []entity.Task{}
	for _, t := range templates {
		next := mapNextScheds[t.TaskScheduleID].NextSchedule
		var endPtr *time.Time
		if mapNextScheds[t.TaskScheduleID].Duration > 0 {
			end := next.Add(time.Duration(mapNextScheds[t.TaskScheduleID].Duration) * time.Second)
			endPtr = &end
		}
		t.StartDate = &next
		t.EndDate = endPtr
		t.IsTemplate = false
		tasks = append(tasks, t)
	}

	_, err = a.d.TaskUC.Create(ctxx.New(ctx, entity.User{}), nil, tasks)
	if err != nil {
		return err
	}

	return nil
}

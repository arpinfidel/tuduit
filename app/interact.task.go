package app

import (
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
)

type TaskListParams struct {
	Page int `rose:"page,p,default=1"`
	Size int `rose:"size,n,default=25"`

	Search    string `rose:"search,s,q"`
	UserName  string `rose:"username,u"`
	Started   *bool  `rose:"started"`
	Completed *bool  `rose:"completed,default=false"`
	Archived  *bool  `rose:"archived,default=false"`
	// Tags     []string `rose:"tags,t"`
}

type TaskListResults struct {
	Tasks []entity.TaskOverview
	Page  string
	Total int
}

func (h *App) GetTaskList(ctx *ctxx.Context, p TaskListParams) (res TaskListResults, err error) {
	userID := ctx.UserID

	if p.UserName != "" {
		user, _, err := h.d.UserUC.Get(ctx, nil, db.Params{
			Where: []db.Where{
				{
					Field: "username",
					Value: p.UserName,
				},
			},
		})
		if err != nil {
			return res, err
		}
		if len(user) == 0 {
			return res, fmt.Errorf("user not found")
		}

		userID = user[0].ID
	}

	where := []db.Where{
		{
			Field: "user_id",
			Op:    db.EqOp,
			Value: userID,
		},
	}

	if p.Search != "" {
		where = append(where, db.Where{
			Op: db.OrOp,
			Value: []db.Where{
				{
					Field: "name",
					Op:    db.LikeOp,
					Value: fmt.Sprintf("%%%s%%", p.Search),
				},
				{
					Field: "description",
					Value: fmt.Sprintf("%%%s%%", p.Search),
				},
			},
		})
	}

	nullOp := db.IsNullOp
	if p.Started != nil {
		if *p.Started {
			nullOp = db.NotNullOp
		}
		where = append(where, db.Where{
			Field: "started_at",
			Op:    nullOp,
		})
	}

	if p.Completed != nil {
		nullOp := db.IsNullOp
		if *p.Completed {
			nullOp = db.NotNullOp
		}
		where = append(where, db.Where{
			Field: "completed_at",
			Op:    nullOp,
		})
	}

	if p.Archived != nil {
		nullOp := db.IsNullOp
		if *p.Archived {
			nullOp = db.NotNullOp
		}
		where = append(where, db.Where{
			Field: "archived_at",
			Op:    nullOp,
		})
	}

	tasks, count, err := h.d.TaskUC.Get(ctx, nil, db.Params{
		Where: where,
		Pagination: &db.Pagination{
			Limit:  p.Size,
			Offset: (p.Page - 1) * p.Size,
		},
		Sort: []db.Sort{
			{
				Field: "priority",
				Asc:   true,
			},
			{
				Field:      "end_date",
				Asc:        true,
				NullsFirst: false,
			},
			{
				Field:      "start_date",
				Asc:        true,
				NullsFirst: true,
			},
			{
				Field: "id",
				Asc:   true,
			},
		},
	})
	if err != nil {
		return res, err
	}

	for _, task := range tasks {
		res.Tasks = append(res.Tasks, task.Overview())
	}

	res.Page = fmt.Sprintf("%d/%d", p.Page, count/p.Size+1)
	res.Total = count

	return res, nil
}

func TaskListToString(res TaskListResults) string {
	resp := ""
	for _, t := range res.Tasks {
		resp += fmt.Sprintf("%d. (P%d) %s\n", t.ID, t.Priority, t.Name)
		if t.StartDate != "" {
			resp += fmt.Sprintf("\tStart: %s\n", t.StartDate)
		}
		if t.EndDate != "" {
			resp += fmt.Sprintf("\tEnd: %s\n", t.EndDate)
		}
	}

	resp += fmt.Sprintf("Total tasks: %d", res.Total)
	resp += fmt.Sprintf("Page: %s", res.Page)

	return resp
}

type CreateTaskParams struct {
	Name        string     `rose:"name,n,required="`
	Priority    int        `rose:"priority,p,default=2"`
	Description string     `rose:"description,d"`
	StartDate   *time.Time `rose:"start_date,dt"`
	EndDate     *time.Time `rose:"end_date,dt"`
	Assignee    string     `rose:"assignee,a"`
}

func (h *App) CreateTask(ctx *ctxx.Context, p CreateTaskParams) (task entity.TaskOverview, err error) {
	userID := ctx.UserID

	if p.Assignee != "" {
		user, _, err := h.d.UserUC.Get(ctx, nil, db.Params{
			Where: []db.Where{
				{
					Field: "username",
					Value: p.Assignee,
				},
			},
		})
		if err != nil {
			return task, err
		}
		if len(user) == 0 {
			return task, fmt.Errorf("user not found")
		}

		userID = user[0].ID
	}

	t, err := h.d.TaskUC.Create(ctx, nil, []entity.Task{
		{
			UserID:      userID,
			Name:        p.Name,
			Priority:    p.Priority,
			Description: p.Description,
			StartDate:   p.StartDate,
			EndDate:     p.EndDate,
		},
	})
	if err != nil {
		return task, err
	}
	task = t[0].Overview()

	return task, nil
}

type UpdateTaskParams struct {
	ID int `rose:"id,i,required="`

	Started   *bool `rose:"started,s"`
	Completed *bool `rose:"completed,c"`
	Archived  *bool `rose:"archived,a"`

	Name        *string `rose:"name,n"`
	Priority    *int    `rose:"priority,p"`
	Description *string `rose:"description,d"`
	Assignee    *string `rose:"assignee,ass"`
}

func (h *App) UpdateTask(ctx *ctxx.Context, p UpdateTaskParams) (taskO entity.TaskOverview, err error) {
	t, _, err := h.d.TaskUC.GetByIDs(ctx, nil, []int{p.ID}, entity.Pagination{PageSize: 1, Page: 1})
	if err != nil {
		return taskO, err
	}
	if len(t) == 0 {
		return taskO, fmt.Errorf("task not found")
	}

	now := time.Now()

	if p.Started != nil {
		if *p.Started {
			t[0].StartedAt = &now
		} else {
			t[0].StartedAt = nil
		}
	}

	if p.Completed != nil {
		if *p.Completed {
			t[0].CompletedAt = &now
		} else {
			t[0].CompletedAt = nil
		}
	}

	if p.Archived != nil {
		if *p.Archived {
			t[0].ArchivedAt = &now
		} else {
			t[0].ArchivedAt = nil
		}
	}

	if p.Name != nil {
		t[0].Name = *p.Name
	}

	if p.Priority != nil {
		t[0].Priority = *p.Priority
	}

	if p.Description != nil {
		t[0].Description = *p.Description
	}

	if p.Assignee != nil {
		user, _, err := h.d.UserUC.Get(ctx, nil, db.Params{
			Where: []db.Where{
				{
					Field: "username",
					Value: *p.Assignee,
				},
			},
		})
		if err != nil {
			return taskO, err
		}
		if len(user) == 0 {
			return taskO, fmt.Errorf("user not found")
		}
		t[0].UserID = user[0].ID
	}

	task, err := h.d.TaskUC.Update(ctx, nil, t[0])
	if err != nil {
		return taskO, err
	}

	return task.Overview(), nil
}

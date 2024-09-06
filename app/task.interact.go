package app

import (
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/rose"
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
	Tasks []entity.TaskOverview `rose:"tasks"`
	Page  string                `rose:"page"`
	Total int                   `rose:"total"`
}

func (a *App) GetTaskList(ctx *ctxx.Context, p TaskListParams) (res TaskListResults, err error) {
	userID := ctx.User.ID

	if p.UserName != "" {
		user, _, err := a.d.UserUC.Get(ctx, nil, db.Params{
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
			Field: "user_ids",
			Op:    db.ArrContains,
			Value: []int64{userID},
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

	tasks, count, err := a.d.TaskUC.Get(ctx, nil, db.Params{
		WithCount: true,
		Where:     where,
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

	rose.ChangeTimezone(&tasks, ctx.User.Timezone())

	for _, task := range tasks {
		res.Tasks = append(res.Tasks, task.Overview())
	}

	res.Page = fmt.Sprintf("%d/%d", p.Page, count/p.Size+1)
	res.Total = count

	return res, nil
}

func TaskListToString(ctx *ctxx.Context, res TaskListResults) string {
	resp := ""
	for i, t := range res.Tasks {
		if i > 0 && t.Priority > res.Tasks[i-1].Priority {
			resp += "\n"
		}

		resp += fmt.Sprintf("%s. (P%d) %s\n", t.ID, t.Priority, t.Name)
		if t.StartDate != "" {
			resp += fmt.Sprintf("\tStart: %s\n", t.StartDate)
		}
		if t.EndDate != "" {
			resp += fmt.Sprintf("\tEnd: %s\n", t.EndDate)
		}
	}

	resp += fmt.Sprintf("Total tasks: %d\n", res.Total)
	resp += fmt.Sprintf("Page: %s\n", res.Page)

	return resp
}

type CreateTaskParams struct {
	Name        string     `rose:"name,n,required="`
	Priority    int        `rose:"priority,p,default=2"`
	Description string     `rose:"description,d"`
	StartDate   *time.Time `rose:"start_date,sd"`
	EndDate     *time.Time `rose:"end_date,ed"`
	Assignees   []string   `rose:"assignees,a"`
}

func (a *App) CreateTask(ctx *ctxx.Context, p CreateTaskParams) (task entity.TaskOverview, err error) {
	userIDs := []int64{ctx.User.ID}

	if len(p.Assignees) > 0 {
		user, _, err := a.d.UserUC.Get(ctx, nil, db.Params{
			Where: []db.Where{
				{
					Field: "username",
					Value: p.Assignees,
				},
			},
		})
		if err != nil {
			return task, err
		}
		if len(user) != len(p.Assignees) {
			return task, fmt.Errorf("user not found")
		}

		userIDs = []int64{}
		for _, u := range user {
			userIDs = append(userIDs, u.ID)
		}
	}

	t, err := a.d.TaskUC.Create(ctx, nil, []entity.Task{
		{
			UserIDs:     userIDs,
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
	IDs []entity.Base36[uint64] `rose:"ids,i,required="`

	Started   *bool `rose:"started,s"`
	Completed *bool `rose:"completed,c"`
	Archived  *bool `rose:"archived,a"`

	Name        *string  `rose:"name,n"`
	Priority    *int     `rose:"priority,p"`
	Description *string  `rose:"description,d"`
	Assignees   []string `rose:"assignees,ass"`
}

func (a *App) UpdateTask(ctx *ctxx.Context, p UpdateTaskParams) (taskOs []entity.TaskOverview, err error) {
	ids := []int64{}
	for _, id := range p.IDs {
		ids = append(ids, int64(id.V))
	}
	tasks, _, err := a.d.TaskUC.GetByIDs(ctx, nil, ids, entity.Pagination{PageSize: len(p.IDs), Page: 1})
	if err != nil {
		return taskOs, err
	}
	if len(tasks) < len(p.IDs) {
		return taskOs, fmt.Errorf("task not found")
	}

	now := time.Now()
	taskOs = []entity.TaskOverview{}

	for _, t := range tasks {
		if p.Started != nil {
			if *p.Started {
				t.StartedAt = &now
			} else {
				t.StartedAt = nil
			}
		}

		if p.Completed != nil {
			if *p.Completed {
				t.CompletedAt = &now
			} else {
				t.CompletedAt = nil
			}
		}

		if p.Archived != nil {
			if *p.Archived {
				t.ArchivedAt = &now
			} else {
				t.ArchivedAt = nil
			}
		}

		if p.Name != nil {
			t.Name = *p.Name
		}

		if p.Priority != nil {
			t.Priority = *p.Priority
		}

		if p.Description != nil {
			t.Description = *p.Description
		}

		if len(p.Assignees) > 0 {
			user, _, err := a.d.UserUC.Get(ctx, nil, db.Params{
				Where: []db.Where{
					{
						Field: "username",
						Value: p.Assignees,
					},
				},
			})
			if err != nil {
				return taskOs, err
			}
			if len(user) != len(p.Assignees) {
				return taskOs, fmt.Errorf("user not found")
			}

			t.UserIDs = []int64{}
			for _, u := range user {
				t.UserIDs = append(t.UserIDs, u.ID)
			}
		}

		task, err := a.d.TaskUC.Update(ctx, nil, t)
		if err != nil {
			return taskOs, err
		}

		taskOs = append(taskOs, task.Overview())
	}

	return taskOs, nil
}

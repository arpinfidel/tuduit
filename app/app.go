package app

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	taskuc "github.com/arpinfidel/tuduit/usecase/task"
	useruc "github.com/arpinfidel/tuduit/usecase/user"
)

type Context struct {
	context.Context

	UserID int64
}

type App struct {
	d   Dependencies
	cfg Config
}

type Dependencies struct {
	TaskUC *taskuc.UseCase
	UserUC *useruc.UseCase
}

type Config struct{}

type TaskListParams struct {
	UserName string `rose:"username,u"`
	Page     int    `rose:"page,p"`
	Size     int    `rose:"size,s,n"`
}

func (h *App) TaskList(ctx *Context, p TaskListParams) (tasks []entity.Task, cnt int, err error) {
	var userID int64 = -1
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
			return nil, 0, err
		}
		userID = user[0].ID
	}

	tasks, cnt, err = h.d.TaskUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "user_id",
				Value: userID,
			},
		},
		Pagination: &db.Pagination{
			Limit:  p.Size,
			Offset: (p.Page - 1) * p.Size,
		},
	})
	if err != nil {
		return nil, 0, err
	}

	return tasks, cnt, nil
}

package business

import (
	"github.com/arpinfidel/tuduit/app"
	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	taskuc "github.com/arpinfidel/tuduit/usecase/task"
	useruc "github.com/arpinfidel/tuduit/usecase/user"
)

type Business struct {
	d   Dependencies
	cfg Config
}

type Dependencies struct {
	TaskUC *taskuc.UseCase
	UserUC *useruc.UseCase
}

type Config struct{}

func (h *Business) List(ctx *app.Context, username string, page, size int) (tasks []entity.Task, cnt int, err error) {
	var userID int64 = -1
	if username != "" {
		user, _, err := h.d.UserUC.Get(ctx, nil, db.Params{
			Where: []db.Where{
				{
					Field: "username",
					Value: username,
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
			Limit:  size,
			Offset: (page - 1) * size,
		},
	})
	if err != nil {
		return nil, 0, err
	}

	return tasks, cnt, nil
}

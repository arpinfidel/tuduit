package app

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
)

func (a *App) GetUserByIDs(ctx context.Context, ids []int64, pg entity.Pagination) (data []entity.User, total int, err error) {
	u, total, err := a.d.UserUC.GetByIDs(ctx, nil, ids, pg)
	return u, total, err
}

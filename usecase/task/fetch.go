package taskuc

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/jmoiron/sqlx"
)

func (u *UseCase) Get(ctx context.Context, dbTx *sqlx.Tx, param db.Params) (data []entity.Task, total int, err error) {
	param.Where = append(param.Where, db.Where{
		Field: "is_template",
		Op:    db.EqOp,
		Value: false,
	})
	return u.IRepo.Get(ctx, dbTx, param)
}

func (u *UseCase) GetTemplate(ctx context.Context, dbTx *sqlx.Tx, param db.Params) (data []entity.Task, total int, err error) {
	param.Where = append(param.Where, db.Where{
		Field: "is_template",
		Op:    db.EqOp,
		Value: true,
	})
	return u.IRepo.Get(ctx, dbTx, param)
}

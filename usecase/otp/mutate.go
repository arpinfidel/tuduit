package otpuc

import (
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/jmoiron/sqlx"
)

func (u *UseCase) Create(ctx *ctxx.Context, dbTx *sqlx.Tx, newData []entity.OTP) (data []entity.OTP, err error) {
	for i := range newData {
		d := &newData[i]
		d.CreatedAt = time.Now()
		d.CreatedBy = ctx.User.ID
		d.UpdatedAt = time.Now()
		d.UpdatedBy = ctx.User.ID
	}

	return u.IRepo.Create(ctx, dbTx, newData)
}

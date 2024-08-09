package taskuc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/jmoiron/sqlx"
)

func (u *UseCase) Create(ctx *ctxx.Context, dbTx *sqlx.Tx, newData []entity.Task) (data []entity.Task, err error) {
	fmt.Printf(" >> debug >> ctx.UserID: %#v\n", ctx.UserID)
	for i := range newData {
		d := &newData[i]
		d.CreatedAt = time.Now()
		d.CreatedBy = ctx.UserID
		d.UpdatedAt = time.Now()
		d.UpdatedBy = ctx.UserID
	}

	fmt.Printf(" >> debug >> newData: %s\n", func() string { j, _ := json.MarshalIndent(newData, "", "  "); return string(j) }())

	return u.IRepo.Create(ctx, dbTx, newData)
}

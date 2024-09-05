package app

import (
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/pkg/errs"
)

var _ = registerTask("CreatePartitions", "* * * * *", func() func() error { return a.CreatePartitions })

func (a *App) CreatePartitions() (err error) {
	defer errs.DeferTrace(&err)()

	monthly := []string{
		"trx_otp",
	}

	now := time.Now().UTC()

	trx, err := a.d.DB.Master.Begin()
	if err != nil {
		return err
	}
	defer trx.Rollback()

	for _, table := range monthly {
		a.l.Infof("Creating partitions for %s", table)
		for i := 0; i < 3; i++ {
			futureDate := now.AddDate(0, i, 0)
			suffix := futureDate.Format("2006_01") // yyyy_mm
			part := fmt.Sprintf("%s_%s", table, suffix)
			start := time.Date(futureDate.Year(), futureDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(futureDate.Year(), futureDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)

			q := `
				create table if not exists %s partition of %s
					for values from ('%s') to ('%s');
			`
			q = fmt.Sprintf(q, part, table, start.Format("2006-01-02"), end.Format("2006-01-02"))

			a.l.Infof("  Creating partition %s", part)
			_, err = trx.Exec(q)
			if err != nil {
				return err
			}
		}
	}

	err = trx.Commit()
	if err != nil {
		return err
	}

	return nil
}

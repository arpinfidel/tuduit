package app

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/cron"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/log"
	checkinuc "github.com/arpinfidel/tuduit/usecase/checkin"
	taskuc "github.com/arpinfidel/tuduit/usecase/task"
	useruc "github.com/arpinfidel/tuduit/usecase/user"
	"github.com/jmoiron/sqlx"
	"go.mau.fi/whatsmeow"
)

type App struct {
	l *log.Logger

	d   Dependencies
	cfg Config
}

type Dependencies struct {
	TaskUC    *taskuc.UseCase
	UserUC    *useruc.UseCase
	CheckinUC *checkinuc.UseCase

	Cron     *cron.Cron
	WaClient *whatsmeow.Client
}

type Config struct{}

var a *App

func New(logger *log.Logger, deps Dependencies, cfg Config) *App {
	if a != nil {
		return a
	}

	a = &App{
		l: logger,

		d:   deps,
		cfg: cfg,
	}

	for _, f := range tasks {
		a.d.Cron.RegisterJob(f.Name, f.Schedule, f.Func)
	}

	err := a.SendCheckInMsgs()
	if err != nil {
		a.l.Errorf("SendCheckInMsgs: %v", err)
	}

	return a
}

type job struct {
	Name     string
	Schedule string
	Func     func() error
}

var (
	tasks = []job{}
)

func registerTask(name, schedule string, f func() error) struct{} {
	tasks = append(tasks, job{
		Name:     name,
		Schedule: schedule,
		Func:     f,
	})
	return struct{}{}
}

func (a *App) GetUser(ctx context.Context, dbTx *sqlx.Tx, param db.Params) (data []entity.User, total int, err error) {
	return a.d.UserUC.Get(ctx, dbTx, param)
}

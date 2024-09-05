package app

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/cron"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/jwt"
	"github.com/arpinfidel/tuduit/pkg/log"
	checkinuc "github.com/arpinfidel/tuduit/usecase/checkin"
	otpuc "github.com/arpinfidel/tuduit/usecase/otp"
	scheduleuc "github.com/arpinfidel/tuduit/usecase/schedule"
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
	TaskUC     *taskuc.UseCase
	ScheduleUC *scheduleuc.UseCase
	UserUC     *useruc.UseCase
	CheckinUC  *checkinuc.UseCase
	OTPUC      *otpuc.UseCase

	Cron *cron.Cron
	JWT  *jwt.JWT

	WaClient *whatsmeow.Client

	DB *db.DB
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

	a.l.Infof("Registering tasks: %d", len(tasks))
	for _, f := range tasks {
		a.l.Infof("Registering task: %s at %s", f.Name, f.Schedule)
		a.d.Cron.RegisterJob(f.Schedule, f.Name, f.Func())
	}

	a.l.Infof("Creating partitions")
	err := a.CreatePartitions()
	if err != nil {
		a.l.Errorf("CreatePartitions: %v", err)
	}
	a.l.Infof("Partitions created")

	err = a.d.Cron.Start()
	if err != nil {
		a.l.Errorf("Cron.Start: %v", err)
	}

	return a
}

type job struct {
	Name     string
	Schedule string
	Func     func() func() error
}

var (
	tasks = []job{}
)

func registerTask(name, schedule string, f func() func() error) struct{} {
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

package cron

import (
	"context"

	"github.com/arpinfidel/tuduit/pkg/log"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	ctx context.Context
	l   *log.Logger

	client *cron.Cron
}

func New(ctx context.Context, l *log.Logger) *Cron {
	return &Cron{
		ctx: ctx,
		l:   l,

		client: cron.New(),
	}
}

type Job struct {
	Name string
	Func func() error
}

func (c *Cron) RecoverPanic(f Job) Job {
	return Job{
		Name: f.Name,
		Func: func() error {
			defer func() {
				if err := recover(); err != nil {
					c.l.Errorf("panic: %v", err)
				}
			}()
			return f.Func()
		},
	}
}

func (c *Cron) LogStartEnd(f Job) Job {
	return Job{
		Name: f.Name,
		Func: func() error {
			c.l.Infof("start %s", f.Name)
			defer c.l.Infof("end %s", f.Name)
			return f.Func()
		},
	}
}

func (c *Cron) Wrap(f Job) func() {
	f = c.LogStartEnd(f)
	f = c.RecoverPanic(f)
	return func() {
		if err := f.Func(); err != nil {
			c.l.Errorf(err.Error())
		}
	}
}

func (c *Cron) RegisterJob(schedule string, name string, f func() error) {
	_, err := c.client.AddFunc(schedule, c.Wrap(Job{
		Name: name,
		Func: f,
	}))
	if err != nil {
		panic(err)
	}
}

func (c *Cron) Start() error {
	c.client.Start()
	return nil
}

package handler

import (
	"fmt"

	"github.com/arpinfidel/tuduit/app"
	taskuc "github.com/arpinfidel/tuduit/usecase/task"
	useruc "github.com/arpinfidel/tuduit/usecase/user"
	"github.com/urfave/cli/v2"
)

type Handler struct {
	d   Dependencies
	cfg Config

	Events Events
}

type Dependencies struct {
	App *app.App

	TaskUC *taskuc.UseCase
	UserUC *useruc.UseCase
}

type Config struct {
	OutputType string
}

type Events struct {
	Output chan string
}

func New(deps Dependencies, cfg Config) *Handler {
	return &Handler{
		d:   deps,
		cfg: cfg,

		Events: Events{
			Output: make(chan string, 1000),
		},
	}
}

type ActionFunc func(ctx *app.Context, c *cli.Context) (err error)

func (h *Handler) output(output string) (err error) {
	switch h.cfg.OutputType {
	case "string":
		h.Events.Output <- output
	case "stdout":
		fmt.Println(output)
	}
	return nil
}

func (h *Handler) List() (flags []cli.Flag, actionFunc ActionFunc) {
	type flagsType struct {
		Size     int    `tuduit:"size"`
		Username string `tuduit:"username"`
		Search   string `tuduit:"search"`
	}

	flags, err := h.makeFlags(flagsType{})
	if err != nil {
		return nil, nil
	}

	actionFunc = func(ctx *app.Context, c *cli.Context) (err error) {
		args := struct {
			Page int `tuduit:"page"`
		}{
			Page: 1,
		}

		_, err = h.getArgs(ctx, c, &args)
		if err != nil {
			return err
		}

		flags := flagsType{
			Size: 10,
		}
		err = h.getFlags(ctx, c, &flags)
		if err != nil {
			return err
		}

		tasks, cnt, err := h.d.App.TaskList(ctx, app.TaskListParams{
			UserName: flags.Username,
			Page:     args.Page,
			Size:     flags.Size,
		})
		if err != nil {
			return err
		}

		msg := ""
		for _, task := range tasks {
			msg += fmt.Sprintf("%s\n", task.Name)
		}

		msg += fmt.Sprintf("page %d/%d", args.Page, cnt/flags.Size+1)

		return h.output(msg)
	}

	return flags, actionFunc
}

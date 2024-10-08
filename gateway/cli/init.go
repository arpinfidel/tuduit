package cli

import (
	"github.com/arpinfidel/tuduit/gateway/cli/handler"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/urfave/cli/v2"
)

type App struct {
	ctx *ctxx.Context

	app cli.App

	handler *handler.Handler
}

type Dependencies struct {
	handler.Dependencies
}

type Config struct {
	handler.Config
}

func New(ctx *ctxx.Context, deps Dependencies, cfg Config) *App {
	cliApp := &cli.App{
		Name:  "todo",
		Usage: "Todo CLI",
	}

	app := &App{
		ctx: ctx,

		app: *cliApp,

		handler: handler.New(deps.Dependencies, cfg.Config),
	}

	app.registerSubCommands()

	return app
}

func command(ctx *ctxx.Context, com *cli.Command, f func() (flags []cli.Flag, actionFunc handler.ActionFunc)) *cli.Command {
	flags, actionFunc := f()
	com.Flags = flags
	com.Action = func(c *cli.Context) error {
		return actionFunc(ctx, c)
	}
	return com
}

func (a *App) registerSubCommands() {
	a.app.Commands = append(a.app.Commands, command(a.ctx, &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List tasks",
	}, a.handler.List))
}

package wabot

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arpinfidel/tuduit/app"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/errs"
	"github.com/arpinfidel/tuduit/pkg/log"
	"github.com/arpinfidel/tuduit/pkg/messenger"
	"github.com/arpinfidel/tuduit/pkg/rose"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

const (
	prefix = "."
)

type WaBot struct {
	ctx context.Context
	l   *log.Logger

	d Dependencies
}

type Dependencies struct {
	Messenger messenger.Messenger

	App *app.App
}

func New(ctx context.Context, l *log.Logger, deps Dependencies) *WaBot {
	s := &WaBot{
		ctx: ctx,
		l:   l,

		d: deps,
	}
	return s
}

func wrapHandler[T any, U any](f func(ctx *ctxx.Context, req T) (resp U, err error)) func(ctx *ctxx.Context, text string) (string, error) {
	type help[T any] struct {
		T    T    `rose:"flatten="`
		Help bool `rose:"help,h"`
	}
	return func(ctx *ctxx.Context, text string) (string, error) {
		req := help[T]{}
		r, err := rose.NewParser(ctx, prefix).ParseTextMsg(text, &req)
		if err != nil {
			return "internal error", err
		}

		if req.Help {
			return rose.Help(req.T)
		}

		if len(r.Errors) > 0 {
			errStrings := []string{}
			for _, e := range r.Errors {
				errStrings = append(errStrings, e.Error())
			}
			return "", errors.New("\n" + strings.Join(errStrings, "\n"))
		}

		if !r.Valid {
			resp, err := rose.Help(req.T)
			if err != nil {
				return "", err
			}
			return "invalid request: " + r.Errors[0].Error() + "\n" + resp, nil
		}

		resp, err := f(ctx, req.T)
		if err != nil {
			return "", err
		}

		rose.ChangeTimezone(&resp, ctx.User.Timezone())

		respStr, err := yaml.Marshal(resp)
		if err != nil {
			return "", err
		}

		return "```\n" + string(respStr) + "\n```", nil
	}
}

func (s *WaBot) eventHandler(ctx context.Context, msg *messenger.Message) error {
	usrs, _, err := s.d.App.GetUser(s.ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "whatsapp_number",
				Value: msg.SenderID,
			},
		},
	})
	if err != nil {
		s.l.Errorf("failed to get user: %v", err)
		if trace := errs.GetTrace(err); len(trace) > 0 {
			s.l.Errorln(strings.Join(trace, "\n"))
		}
		return err
	}
	if len(usrs) == 0 {
		s.l.Debugf("user not registered: %s", msg.SenderID)
		return nil
	}

	usr := usrs[0]
	tctx := ctxx.New(ctx, usr)
	tctx = ctxx.WithMessage(tctx, msg)

	err = s.handleText(tctx, msg)
	if err != nil {
		s.l.Errorf("failed to handle text message: %v", err)
		return err
	}

	return nil
}

func (s *WaBot) handleText(ctx *ctxx.Context, v *messenger.Message) (err error) {
	var textBlock *messenger.TextBlock

	for _, b := range v.Blocks {
		if tb, ok := b.(*messenger.TextBlock); ok {
			textBlock = tb
			break
		}
	}

	if textBlock == nil {
		return nil
	}

	text := strings.TrimSpace(textBlock.Text)

	if text == "" {
		s.l.Errorf("missing message text: %#v", v)
		return errors.New("missing message text")
	}

	if !strings.HasPrefix(text, prefix+"tuduit") {
		return nil
	}

	text = strings.TrimSpace(text[len(prefix+"tuduit"):])
	parts := strings.SplitN(text, "\n", 2)

	parts2 := strings.SplitN(parts[0], " ", 2)

	command := parts2[0]
	value := ""
	if len(parts2) > 1 {
		value += parts2[1]
	}
	if len(parts) > 1 {
		value += "\n" + parts[1]
	}

	resp, err := s.routeText(ctx, command, value)
	if err != nil {
		s.l.Errorf("error in %v command: %v", command, err.Error())
		msg := messenger.Message{
			Conversation: v.Conversation,
			Blocks: []messenger.Block{
				&messenger.TextBlock{
					Text: fmt.Sprintf("error in %v command: %v", command, err.Error()),
				},
			},
		}
		s.d.Messenger.SendMessage(s.ctx, msg)
		return nil
	}

	msg := messenger.Message{
		Conversation: v.Conversation,
		Blocks: []messenger.Block{
			&messenger.TextBlock{
				Text: resp,
			},
		},
	}

	s.d.Messenger.SendMessage(s.ctx, msg)

	return nil
}

func (s *WaBot) routeText(ctx *ctxx.Context, command string, value string) (resp string, err error) {
	type comm struct {
		group string
		names []string
		f     func(*ctxx.Context, string) (string, error)
	}
	funcs := []comm{
		{
			group: "task",
			names: []string{"task-list", "tl"},
			f:     wrapHandler(s.HandlerTaskList),
		},
		{
			group: "task",
			names: []string{"task-create", "tc"},
			f:     wrapHandler(s.d.App.CreateTask),
		},
		{
			group: "task",
			names: []string{"task-update", "tu"},
			f:     wrapHandler(s.d.App.UpdateTask),
		},
		{
			group: "schedule",
			names: []string{"schedule-create", "sc"},
			f:     wrapHandler(s.d.App.CreateSchedule),
		},
		{
			group: "misc",
			names: []string{"reminder-create", "rc"},
			f:     wrapHandler(s.d.App.SetReminder),
		},
	}

	funcs = append(funcs, comm{
		group: "misc",
		names: []string{"help", "h", ""},
		f: func(ctx *ctxx.Context, s string) (string, error) {
			ss := "available commands:\n"

			groups := map[string][]comm{}
			for _, fs := range funcs {
				groups[fs.group] = append(groups[fs.group], fs)
			}

			for group, fs := range groups {
				ss += fmt.Sprintf("%s:\n", group)

				for _, f := range fs {
					ss += fmt.Sprintf(" - %s\n", strings.Join(f.names, " | "))
				}
			}

			return ss, nil
		},
	})

	fmap := map[string]func(*ctxx.Context, string) (string, error){}
	for _, fs := range funcs {
		for _, f := range fs.names {
			fmap[f] = fs.f
		}
	}

	f, ok := fmap[command]
	if !ok {
		return "", fmt.Errorf("unsupported command: %s", command)
	}

	return f(ctx, value)
}

func (s *WaBot) Start() (err error) {
	s.d.Messenger.AddHandler(context.Background(), s.eventHandler)

	return nil
}

func (s *WaBot) v1() {
}

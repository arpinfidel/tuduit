package wabot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/arpinfidel/tuduit/app"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/log"
	"github.com/arpinfidel/tuduit/pkg/rose"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
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
	WaClient *whatsmeow.Client

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
		Help bool `rose:"help"`
	}
	return func(ctx *ctxx.Context, text string) (string, error) {
		req := help[T]{}
		r, err := rose.NewParser(prefix).ParseTextMsg(text, &req)
		if err != nil {
			return "", err
		}

		if req.Help {
			return rose.Help(req.T)
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

		respStr, err := yaml.Marshal(resp)
		if err != nil {
			return "", err
		}

		return string(respStr), nil
	}
}

func (s *WaBot) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		sender := v.Info.Sender.User

		usrs, _, err := s.d.App.GetUser(s.ctx, nil, db.Params{
			Where: []db.Where{
				{
					Field: "whatsapp_number",
					Value: sender,
				},
			},
		})
		if err != nil {
			s.l.Errorf("failed to get user: %v", err)
			return
		}

		if len(usrs) == 0 {
			s.l.Errorf("user not registered: %s", sender)
			return
		}

		usr := usrs[0]
		ctx := ctxx.New(s.ctx, usr.ID)

		switch v.Info.Type {
		default:
			s.l.Errorf("unsupported message text")
			return
		case "text":
			err = s.handleText(ctx, v)
			if err != nil {
				s.l.Errorf("failed to handle text message: %v", err)
				return
			}
		case "media":
			switch v.Info.MediaType {
			case "image":
				// TODO: implement image handling
			}
		}
	}
}

func (s *WaBot) handleText(ctx *ctxx.Context, v *events.Message) (err error) {
	fmt.Printf(" >> debug >> v.Message: %s\n", func() string { j, _ := json.MarshalIndent(v.Message, "", "  "); return string(j) }())
	fmt.Printf(" >> debug >> v.RawMessage: %s\n", func() string { j, _ := json.MarshalIndent(v.RawMessage, "", "  "); return string(j) }())

	text := ""

	if v.Message != nil && v.Message.ExtendedTextMessage != nil && v.Message.ExtendedTextMessage.Text != nil {
		text = *v.Message.ExtendedTextMessage.Text
	}

	if v.Message != nil && v.Message.Conversation != nil {
		text = *v.Message.Conversation
	}

	if text == "" {
		s.l.Errorf("missing message text: %s", v.Info.Type)
		return errors.New("missing message text")
	}

	if !strings.HasPrefix(text, prefix+"tuduit") {
		return nil
	}

	fmt.Printf(" >> debug >> text: %s\n", text)

	text = strings.TrimSpace(text[len(prefix+"tuduit"):])
	parts := strings.SplitN(text, "\n", 2)
	fmt.Printf(" >> debug >> parts: %s\n", parts)
	parts2 := strings.SplitN(parts[0], " ", 2)
	fmt.Printf(" >> debug >> parts2: %s\n", parts2)
	command := parts2[0]
	value := ""
	if len(parts2) > 1 {
		value += parts2[1]
	}
	if len(parts) > 1 {
		value += "\n" + parts[1]
	}

	fmt.Printf(" >> debug >> value: %s\n", value)

	resp, err := s.routeText(ctx, command, value)
	if err != nil {
		s.l.Errorf("error in %v command: %v", command, err.Error())
		s.d.WaClient.SendMessage(s.ctx, v.Info.Chat, &waE2E.Message{
			Conversation: proto.String(fmt.Sprintf("error in %v command: %v", command, err.Error())),
		})
		return nil
	}

	s.d.WaClient.SendMessage(s.ctx, v.Info.Chat, &waE2E.Message{
		Conversation: proto.String(resp),
	})

	return nil
}

func (s *WaBot) routeText(ctx *ctxx.Context, command string, value string) (resp string, err error) {
	funcs := []struct {
		names []string
		f     func(*ctxx.Context, string) (string, error)
	}{
		{
			names: []string{"list", "l"},
			f:     wrapHandler(s.d.App.GetTaskList),
		},
		{
			names: []string{"create", "c"},
			f:     wrapHandler(s.d.App.CreateTask),
		},
		{
			names: []string{"update", "u"},
			f:     wrapHandler(s.d.App.UpdateTask),
		},
	}

	funcs = append(funcs, struct {
		names []string
		f     func(*ctxx.Context, string) (string, error)
	}{
		names: []string{"help", "h", ""},
		f: func(ctx *ctxx.Context, s string) (string, error) {
			ss := "available commands:\n"

			for _, fs := range funcs {
				ss += fmt.Sprintf(" - %s\n", strings.Join(fs.names, " | "))
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
	s.d.WaClient.AddEventHandler(s.eventHandler)

	return nil
}

func (s *WaBot) v1() {
}

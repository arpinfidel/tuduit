package wabot

import (
	"github.com/arpinfidel/tuduit/app"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
)

func (s *WaBot) HandlerTaskList(ctx *ctxx.Context, req app.TaskListParams) (resp string, err error) {
	res, err := s.d.App.GetTaskList(ctx, req)
	if err != nil {
		return "", err
	}

	resp = app.TaskListToString(res)

	return resp, nil
}

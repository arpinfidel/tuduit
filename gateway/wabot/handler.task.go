package wabot

import (
	"fmt"

	"github.com/arpinfidel/tuduit/app"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
)

func (s *WaBot) HandlerTaskList(ctx *ctxx.Context, req app.TaskListParams) (resp string, err error) {
	res, err := s.d.App.GetTaskList(ctx, req)
	if err != nil {
		return "", err
	}

	resp = ""
	for _, t := range res.Tasks {
		resp += fmt.Sprintf("%d. (P%d) %s\n", t.ID, t.Priority, t.Name)
		if t.StartDate != "" {
			resp += fmt.Sprintf("\tStart: %s\n", t.StartDate)
		}
		if t.EndDate != "" {
			resp += fmt.Sprintf("\tEnd: %s\n", t.EndDate)
		}
	}

	resp += fmt.Sprintf("Total: %d", res.Count)

	return resp, nil
}

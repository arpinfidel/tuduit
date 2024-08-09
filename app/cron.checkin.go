package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"gopkg.in/yaml.v3"
)

var _ = registerTask("SendCheckInMsgs", "* * * * *", a.SendCheckInMsgs)

func (a *App) SendCheckInMsgs() error {
	ctx := context.Background()

	ci, _, err := a.d.CheckinUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Op: db.AndOp,
				Value: []db.Where{
					{
						Field: "time(check_in_time)",
						Op:    db.LtOrEqOp,
						Value: "time('now')",
					},
					{
						Field: "datetime(last_sent)",
						Op:    db.LtOrEqOp,
						Value: "datetime(date('now'), '+'||check_in_time)",
					},
				},
			},
		},
		Pagination: &db.Pagination{
			Limit:  1 << 30,
			Offset: 0,
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf(" >> debug >> ci: %s\n", func() string { j, _ := json.MarshalIndent(ci, "", "  "); return string(j) }())

	userIDsMap := map[int]struct{}{}
	for _, c := range ci {
		userIDsMap[c.UserID] = struct{}{}
	}
	userIDs := []int{}
	for u := range userIDsMap {
		userIDs = append(userIDs, u)
	}

	users, _, err := a.d.UserUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "id",
				Op:    db.InOp,
				Value: userIDs,
			},
		},
	})
	if err != nil {
		return err
	}
	usersMap := map[int]entity.User{}
	for _, u := range users {
		usersMap[u.ID] = u
	}

	for _, c := range ci {
		u, ok := usersMap[c.UserID]
		if !ok {
			return fmt.Errorf("user not found")
		}

		c.LastSent = time.Now()
		_, err = a.d.CheckinUC.Update(ctx, nil, c)
		if err != nil {
			return err
		}

		false_ := false
		list, err := a.GetTaskList(ctxx.New(ctx, 0), TaskListParams{
			Page:      1,
			Size:      10,
			UserName:  u.Username,
			Completed: &false_,
		})
		if err != nil {
			return err
		}

		respB, err := yaml.Marshal(list)
		if err != nil {
			return err
		}
		respStr := string(respB)

		_, err = a.d.WaClient.SendMessage(ctx, types.JID{
			User: u.WhatsappNumber,
		}, &waE2E.Message{
			Conversation: &respStr,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

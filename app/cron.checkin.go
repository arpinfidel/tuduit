package app

import (
	"context"
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

var _ = registerTask("SendCheckInMsgs", "* * * * *", func() func() error { return a.SendCheckInMsgs })

func (a *App) SendCheckInMsgs() error {
	ctx := context.Background()

	now := time.Now().UTC()

	ci, _, err := a.d.CheckinUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Op: db.AndOp,
				Value: []db.Where{
					{
						Field: "check_in_time",
						Op:    db.LtOrEqOp,
						Value: now.Format("15:04:05"),
					},
					{
						Field: "last_sent",
						Op:    db.LtOrEqOp,
						Value: now.Format("2006-01-02"),
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
	if len(ci) == 0 {
		return nil
	}
	a.l.Infof("Found %d checkins", len(ci))

	userIDsSet := map[int64]struct{}{}
	for _, c := range ci {
		userIDsSet[c.UserID] = struct{}{}
	}
	userIDs := []int64{}
	for u := range userIDsSet {
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
	usersMap := map[int64]entity.User{}
	for _, u := range users {
		usersMap[u.ID] = u
	}
	a.l.Infof("Found %d users", len(users))

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

		res, err := a.GetTaskList(ctxx.New(ctx, c.UserID), TaskListParams{})
		if err != nil {
			return err
		}

		resp := TaskListToString(res)

		_, err = a.d.WaClient.SendMessage(ctx, types.NewJID(u.WhatsappNumber, types.DefaultUserServer), &waE2E.Message{
			Conversation: &resp,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

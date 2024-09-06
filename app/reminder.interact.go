package app

import (
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/messenger"
)

type ReminderParams struct {
	Duration time.Duration `rose:"duration,d,default=5m"`
	Name     string        `rose:"name,n"`
}

func (a *App) SetReminder(ctx *ctxx.Context, p ReminderParams) (res string, err error) {
	date := time.Now().Add(p.Duration)
	dateStr := date.Format("2006-01-02 15:04:05")
	elapsed := time.Since(date).Round(time.Second)
	if elapsed < 0 {
		dateStr += fmt.Sprintf(" (in %s)", -elapsed)
	} else {
		dateStr += fmt.Sprintf(" (%s ago)", elapsed)
	}
	msg := fmt.Sprintf("Set reminder for %s", dateStr)

	go func() {
		time.Sleep(p.Duration)

		msg := fmt.Sprintf("Reminder from %s ago", p.Duration)
		if p.Name != "" {
			msg += fmt.Sprintf(": %s", p.Name)
		}
		err = a.d.WaClient.SendMessage(ctx, messenger.Message{
			Conversation: ctx.Message.Conversation,
			Blocks: []messenger.Block{
				&messenger.TextBlock{
					Text: msg,
				},
			},
		})
		if err != nil {
			a.l.Logger.Sugar().Errorf("error sending message %v", err)
		}
	}()

	return msg, nil
}

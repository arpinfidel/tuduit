package whatsapp

import (
	"context"
	"fmt"

	"github.com/arpinfidel/tuduit/pkg/messenger"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type Whatsapp struct {
	d Dependencies
}

type Dependencies struct {
	WaClient *whatsmeow.Client
}

func New(deps Dependencies) *Whatsapp {
	return &Whatsapp{
		d: deps,
	}
}

func makeConversationJID(conversation messenger.Conversation) (types.JID, error) {
	var jid types.JID
	switch conversation.Type {
	default:
		return jid, fmt.Errorf("unknown conversation type: %s", conversation.Type)
	case messenger.ConversationTypeGroup:
		jid = types.NewJID(conversation.GroupID, types.GroupServer)
	case messenger.ConversationTypeUser:
		jid = types.NewJID(conversation.UserID, types.DefaultUserServer)
	}
	return jid, nil
}

func (w *Whatsapp) SendMessage(ctx context.Context, msg messenger.Message) error {
	jid, err := makeConversationJID(msg.Conversation)
	if err != nil {
		return err
	}

	waMsg, err := buildMessageFromBlocks(msg.Blocks)
	if err != nil {
		return err
	}

	_, err = w.d.WaClient.SendMessage(ctx, jid, waMsg)
	return err
}

func (w *Whatsapp) AddHandler(ctx context.Context, handler messenger.Handler) error {
	w.d.WaClient.AddEventHandler(func(evt interface{}) {
		msg := messenger.Message{}
		switch v := evt.(type) {
		case *events.Message:
			msg.SenderID = v.Info.Sender.User

			if v.Info.IsGroup {
				msg.Conversation = messenger.Conversation{
					Type:    messenger.ConversationTypeGroup,
					GroupID: v.Info.Chat.User,
				}
			} else {
				msg.Conversation = messenger.Conversation{
					Type:   messenger.ConversationTypeUser,
					UserID: v.Info.Sender.User,
				}
			}

			if v.Message != nil && v.Message.ExtendedTextMessage != nil && v.Message.ExtendedTextMessage.Text != nil {
				msg.Blocks = append(msg.Blocks, &messenger.TextBlock{
					Text: *v.Message.ExtendedTextMessage.Text,
				})
			}

			if v.Message != nil && v.Message.Conversation != nil {
				msg.Blocks = append(msg.Blocks, &messenger.TextBlock{
					Text: *v.Message.Conversation,
				})
			}

			handler(ctx, &msg)
		}
	})
	return nil
}

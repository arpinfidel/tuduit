package messenger

import "context"

type Messenger interface {
	SendMessage(ctx context.Context, msg Message) error
	AddHandler(ctx context.Context, handler Handler) error
}

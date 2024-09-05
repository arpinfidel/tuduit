package ctxx

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/crypto"
	"go.mau.fi/whatsmeow/types/events"
)

type Context struct {
	context.Context
	RequestID string
	User      entity.User

	WAEvent *events.Message
}

type Key string

func New(ct context.Context, user entity.User) *Context {
	c := &Context{
		RequestID: crypto.RandomString(16),
		Context:   ct,
		User:      user,
	}
	ct = context.WithValue(ct, Key("context"), c)
	c.Context = ct
	return c
}

func Background() *Context {
	return New(context.Background(), entity.User{})
}

func GetContext(ct context.Context) *Context {
	key := Key("context")
	ctx := ct.Value(key)
	if ctx == nil {
		return Background()
	}
	return ctx.(*Context)
}

func WithWhatsAppMessage(ct context.Context, msg *events.Message) *Context {
	c := GetContext(ct)
	c.WAEvent = msg
	return c
}

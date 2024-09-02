package ctxx

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
	"go.mau.fi/whatsmeow/types/events"
)

type Context struct {
	context.Context
	User entity.User

	WAEvent *events.Message
}

type BodyType string

const (
	BodyTypeTextMsg BodyType = "text_msg"
)

type Body struct {
	Type BodyType
	Text string
}

type Key string

func New(ct context.Context, user entity.User) *Context {
	c := &Context{
		Context: ct,
		User:    user,
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
	return ct.Value(key).(*Context)
}

func WithWhatsappMessage(ct context.Context, msg *events.Message) *Context {
	c := GetContext(ct)
	c.WAEvent = msg
	return c
}

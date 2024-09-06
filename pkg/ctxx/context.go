package ctxx

import (
	"context"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/crypto"
	"github.com/arpinfidel/tuduit/pkg/messenger"
)

type Context struct {
	context.Context
	RequestID string
	User      entity.User

	Message *messenger.Message
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

func WithMessage(ct context.Context, msg *messenger.Message) *Context {
	c := GetContext(ct)
	c.Message = msg
	return c
}

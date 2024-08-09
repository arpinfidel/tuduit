package ctxx

import (
	"context"

	"github.com/arpinfidel/tuduit/pkg/rose"
)

type Context struct {
	context.Context

	UserID int
	Body   Body

	response chan any
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

func New(ct context.Context, userID int) *Context {
	c := &Context{
		Context: ct,
		UserID:  userID,

		response: make(chan any, 2),
	}
	ct = context.WithValue(ct, Key("context"), c)
	c.Context = ct
	return c
}

// GetContext returns the user id from the context
func GetContext(ct context.Context) *Context {
	return ct.Value(Key("context")).(*Context)
}

func (c *Context) SetBody(bodyType BodyType, body string) {
	c.Body = Body{
		Type: bodyType,
		Text: body,
	}
}

func (c *Context) Bind(target any) (rose.Rose, error) {
	parser := rose.NewParser(".")
	switch c.Body.Type {
	default:
		return rose.Rose{}, nil
	case BodyTypeTextMsg:
		return parser.ParseTextMsg(c.Body.Text, target)
	}
}

func (c *Context) AwaitResponse() <-chan any {
	return c.response
}

func (c *Context) Respond(resp any) {
	c.response <- resp
}

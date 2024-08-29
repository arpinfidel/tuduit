package ctxx

import (
	"context"

	"github.com/arpinfidel/tuduit/pkg/rose"
	"go.mau.fi/whatsmeow/types/events"
)

type Context struct {
	context.Context

	UserID int64
	Body   Body

	WAEvent *events.Message

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

func New(ct context.Context, userID int64) *Context {
	c := &Context{
		Context: ct,
		UserID:  userID,

		response: make(chan any, 2),
	}
	ct = context.WithValue(ct, Key("context"), c)
	c.Context = ct
	return c
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

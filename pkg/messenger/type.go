package messenger

import "context"

type Handler func(context.Context, *Message) error

type ConversationType string

const (
	ConversationTypeGroup = ConversationType("group")
	ConversationTypeUser  = ConversationType("user")
)

type Conversation struct {
	Type ConversationType

	GroupID string
	UserID  string
}

type Message struct {
	Conversation Conversation
	SenderID     string

	Blocks []Block
}

package whatsapp

import (
	"fmt"

	"github.com/arpinfidel/tuduit/pkg/messenger"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

func buildTextBlock(block *messenger.TextBlock) (string, error) {
	text := block.Text
	if block.Bold {
		text = fmt.Sprintf("**%s**", text)
	}
	if block.Italic {
		text = fmt.Sprintf("*%s*", text)
	}
	if block.Strike {
		text = fmt.Sprintf("~~%s~~", text)
	}
	if block.Code {
		text = fmt.Sprintf("`%s`", text)
	}
	if block.CodeBlock {
		text = fmt.Sprintf("\n```\n%s\n```\n", text)
	}
	return text, nil
}

func buildMessageFromBlocks(blocks []messenger.Block) (*waE2E.Message, error) {
	msg := waE2E.Message{}

	text := ""
	for _, block := range blocks {
		switch block.Type() {
		default:
			return nil, fmt.Errorf("unknown block type: %s", block.Type())
		case messenger.BlockTypeText:
			t, err := buildTextBlock(block.(*messenger.TextBlock))
			if err != nil {
				return nil, err
			}
			text += t
		}
	}

	msg.Conversation = &text

	return &msg, nil
}

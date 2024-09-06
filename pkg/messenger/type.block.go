package messenger

type BlockType string

type Block interface {
	Type() BlockType
}

var BlockTypeText = BlockType("text")

type TextBlock struct {
	Text string

	// formatting
	Bold      bool
	Italic    bool
	Strike    bool
	Code      bool
	CodeBlock bool
}

func (b *TextBlock) Type() BlockType {
	return BlockTypeText
}

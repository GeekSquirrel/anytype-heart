package apimodel

import (
	"testing"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func textBlock(text string) *model.ChatMessageMessageBlock {
	return &model.ChatMessageMessageBlock{
		Content: &model.ChatMessageMessageBlockContentOfText{
			Text: &model.ChatMessageMessageBlockText{Text: text},
		},
	}
}

func TestChatMessageFromProto_BlockComposedMessage(t *testing.T) {
	msg := &model.ChatMessage{
		Id: "msg-1",
		Blocks: []*model.ChatMessageMessageBlock{
			textBlock("hello bot!"),
			{
				Content: &model.ChatMessageMessageBlockContentOfMessageQuote{
					MessageQuote: &model.ChatMessageMessageBlockMessageQuote{
						MessageId:     "msg-0",
						ParticipantId: "participant-1",
						Content:       &model.ChatMessageMessageBlockText{Text: "quoted"},
					},
				},
			},
			{
				Content: &model.ChatMessageMessageBlockContentOfLink{
					Link: &model.ChatMessageMessageBlockLink{
						TargetObjectId: "obj-1",
						Type:           model.ChatMessageMessageBlockLink_Image,
					},
				},
			},
		},
	}

	cm := ChatMessageFromProto(msg)

	// Desktop clients leave the flat message part empty; blocks carry the text.
	assert.Empty(t, cm.Content.Text)
	assert.Equal(t, []ChatMessageBlock{
		{Type: "text", Text: "hello bot!", Style: "paragraph", Marks: []TextMark{}},
		{Type: "message_quote", MessageId: "msg-0", ParticipantId: "participant-1", Text: "quoted", Marks: []TextMark{}},
		{Type: "link", TargetObjectId: "obj-1", LinkType: "image"},
	}, cm.Blocks)
}

func TestChatMessageFromProto_FlatMessage(t *testing.T) {
	msg := &model.ChatMessage{
		Id: "msg-2",
		Message: &model.ChatMessageMessageContent{
			Text:  "hello from the API",
			Style: model.BlockContentText_Paragraph,
		},
	}

	cm := ChatMessageFromProto(msg)

	assert.Equal(t, "hello from the API", cm.Content.Text)
	assert.Empty(t, cm.Blocks)
}

func TestChatMessageFromProto_NilMessage(t *testing.T) {
	cm := ChatMessageFromProto(nil)
	assert.Empty(t, cm.Content.Text)
	assert.NotNil(t, cm.Blocks)
	assert.NotNil(t, cm.Attachments)
	assert.NotNil(t, cm.Reactions)
}

func TestMessageContentToProto_SynthesizesTextBlock(t *testing.T) {
	msg := MessageContentToProto(AddChatMessageRequest{Text: "from the api"})

	assert.Equal(t, "from the api", msg.Message.Text)
	require.Len(t, msg.Blocks, 1)
	tb := msg.Blocks[0].GetText()
	require.NotNil(t, tb)
	assert.Equal(t, "from the api", tb.Text)
}

func TestEditContentToProto_SynthesizesTextBlock(t *testing.T) {
	msg := EditContentToProto(EditChatMessageRequest{Text: "edited"})

	assert.Equal(t, "edited", msg.Message.Text)
	require.Len(t, msg.Blocks, 1)
	require.NotNil(t, msg.Blocks[0].GetText())
	assert.Equal(t, "edited", msg.Blocks[0].GetText().Text)
}

package slack

import (
	"github.com/slack-go/slack"
)

type MsgOptionsBuilder struct {
	msg              string
	msgParams        slack.PostMessageParameters
	attachmentFields []slack.AttachmentField
}

func NewMsgOptionsBuilder() *MsgOptionsBuilder {
	return &MsgOptionsBuilder{
		msgParams: slack.PostMessageParameters{Markdown: true}, // message is formatted using Slack's markdown-like syntax
	}
}

func (b *MsgOptionsBuilder) SetMessage(msg string) *MsgOptionsBuilder {
	b.msg = msg
	return b
}

func (b *MsgOptionsBuilder) AddAttachmentField(title, value string) *MsgOptionsBuilder {
	b.attachmentFields = append(b.attachmentFields, slack.AttachmentField{
		Title: title,
		Value: value,
		Short: false, // each field is on the new line
	})

	return b
}

func (b *MsgOptionsBuilder) Build() []slack.MsgOption {
	attachment := slack.Attachment{
		Pretext: b.msg,
		Fields:  b.attachmentFields,
	}

	return []slack.MsgOption{
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionPostMessageParameters(b.msgParams),
	}
}

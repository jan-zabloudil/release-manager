package slack

import (
	"fmt"
	"net/url"

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
	return b.addAttachmentField(title, value)
}

func (b *MsgOptionsBuilder) AddAttachmentFieldWithLink(title string, linkURL url.URL, linkText string) *MsgOptionsBuilder {
	// Docs: https://api.slack.com/reference/surfaces/formatting#linking
	link := fmt.Sprintf("<%s|%s>", linkURL.String(), linkText)
	return b.addAttachmentField(title, link)
}

func (b *MsgOptionsBuilder) addAttachmentField(title, value string) *MsgOptionsBuilder {
	// Some of the fields can be empty (e.g. release notes).
	// But we still want to show all fields in the message to keep the consistency and let user know that the value for field was not set.
	if value == "" {
		value = "-"
	}

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

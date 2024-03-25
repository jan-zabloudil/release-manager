package model

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	messageColor = "#36a64f"
	footerText   = "Release Manager"
)

type ReleaseMessage struct {
	Attachments *slack.Attachment
}

func NewReleaseMessage(title, text string) *ReleaseMessage {
	return &ReleaseMessage{
		Attachments: &slack.Attachment{
			Color:   messageColor,
			Pretext: fmt.Sprintf("*%s*\n\n%s", title, text),
			Footer:  footerText,
		},
	}
}

func (m *ReleaseMessage) WithField(title, value string) *ReleaseMessage {
	m.Attachments.Fields = append(m.Attachments.Fields, slack.AttachmentField{
		Title: title,
		Value: value,
		Short: false,
	})

	return m
}

func (m *ReleaseMessage) Options() []slack.MsgOption {
	return []slack.MsgOption{
		slack.MsgOptionAttachments(*m.Attachments),
	}
}

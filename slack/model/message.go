package model

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	messageColor = "#36a64f"
	footerText   = "Release Manager"
)

type Message struct {
	Attachments *slack.Attachment
}

func NewMessage(title, text string) *Message {
	return &Message{
		Attachments: &slack.Attachment{
			Color:   messageColor,
			Pretext: fmt.Sprintf("*%s*\n\n%s", title, text),
			Footer:  footerText,
		},
	}
}

func (m *Message) WithField(title, value string) *Message {
	m.Attachments.Fields = append(m.Attachments.Fields, slack.AttachmentField{
		Title: title,
		Value: value,
		Short: true,
	})

	return m
}

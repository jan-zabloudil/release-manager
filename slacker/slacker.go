package slacker

import (
	"context"
	"log/slog"

	"github.com/slack-go/slack"
)

type Slacker struct {
	client *slack.Client
}

type SilentSlacker struct{}

func New(apiKey string) *Slacker {
	return &Slacker{
		client: slack.New(apiKey),
	}
}

// NewSilent SilentSlacker should be injected into services if no slack token is provided.
// SilentSlacker does not post any messageS to Slack
func NewSilent() *SilentSlacker {
	return &SilentSlacker{}
}

func (s *Slacker) sendTextMessage(ctx context.Context, channelID string, text string) error {
	msg := slack.MsgOptionText(text, false)
	params := slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{Markdown: true})

	if _, _, err := s.client.PostMessageContext(ctx, channelID, msg, params); err != nil {
		return err
	}

	slog.Debug("slack message sent", "channelId", channelID)

	return nil
}

package slack

import (
	"context"
	"log/slog"

	"github.com/slack-go/slack"
)

type Slack struct {
	client *slack.Client
}

type SilentSlack struct{}

// TODO move to interface
func New(apiKey string) *Slack {
	return &Slack{
		client: slack.New(apiKey),
	}
}

// NewSilent SilentSlack should be injected into services if no slack token is provided.
// SilentSlack does not post any messageS to Slack
func NewSilent() *SilentSlack {
	return &SilentSlack{}
}

func (s *Slack) sendTextMessage(ctx context.Context, channelID string, text string) error {
	msg := slack.MsgOptionText(text, false)
	params := slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{Markdown: true})

	if _, _, err := s.client.PostMessageContext(ctx, channelID, msg, params); err != nil {
		return err
	}

	slog.Debug("slack message sent", "channelId", channelID)

	return nil
}

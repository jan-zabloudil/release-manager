package slack

import (
	"context"
	"fmt"
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

func (s *Slack) sendMessage(ctx context.Context, channelID string, options []slack.MsgOption) error {
	_, _, err := s.client.PostMessageContext(
		ctx,
		channelID,
		options...,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	slog.Debug("slack message sent", "channelId", channelID)

	return nil
}

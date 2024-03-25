package slack

import (
	"context"
	"log/slog"

	"release-manager/slack/model"

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

func (s *Slack) sendMessageWithAttachments(ctx context.Context, channelID string, msg *model.Message) error {

	a := msg.Attachments
	_, _, err := s.client.PostMessageContext(
		ctx,
		channelID,
		slack.MsgOptionAttachments(*a),
	)
	if err != nil {
		return err
	}

	slog.Debug("slack message sent", "channelId", channelID)

	return nil
}

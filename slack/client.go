package slack

import (
	"context"
	"fmt"

	svcerrors "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/slack-go/slack"
)

const (
	errInvalidAuth     = "invalid_auth"
	errChannelNotFound = "channel_not_found"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SendReleaseNotification(ctx context.Context, token, channelID string, n model.ReleaseNotification) error {
	msgOptions := NewMsgOptionsBuilder().
		SetMessage(n.Message)

	if n.ProjectName != nil {
		msgOptions.AddAttachmentField("Project", *n.ProjectName)
	}
	if n.ReleaseTitle != nil {
		msgOptions.AddAttachmentField("Release", *n.ReleaseTitle)
	}
	if n.ReleaseNotes != nil {
		msgOptions.AddAttachmentField("Release notes", *n.ReleaseNotes)
	}
	if n.GitTagName != nil && n.GitTagURL != nil {
		msgOptions.AddAttachmentFieldWithLink("Source code", *n.GitTagURL, *n.GitTagName)
	}

	return c.sendMessage(ctx, token, channelID, msgOptions.Build())
}

func (c *Client) sendMessage(ctx context.Context, token, channelID string, msgOptions []slack.MsgOption) error {
	client := slack.New(token)

	if _, _, err := client.PostMessageContext(ctx, channelID, msgOptions...); err != nil {
		switch err.Error() {
		case errInvalidAuth:
			return svcerrors.NewSlackClientUnauthorizedError().Wrap(err)
		case errChannelNotFound:
			return svcerrors.NewSlackChannelNotFoundError().Wrap(err).
				WithMessage(fmt.Sprintf("slack channel (ID: %s) not found", channelID))
		default:
			return fmt.Errorf("failed to send message to slack channel (ID: %s): %w", channelID, err)
		}
	}

	return nil
}

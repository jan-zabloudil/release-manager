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

func (c *Client) SendReleaseNotification(ctx context.Context, tkn model.SlackToken, channelID string, n model.ReleaseNotification) error {
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
	if n.DeployedToEnvironment != nil {
		if n.DeployedServiceURL != nil {
			msgOptions.AddAttachmentFieldWithLink("Deployed to", *n.DeployedServiceURL, *n.DeployedToEnvironment)
		} else {
			msgOptions.AddAttachmentField("Deployed to", *n.DeployedToEnvironment)
		}
	}
	if n.DeployedAt != nil {
		msgOptions.AddAttachmentField("Deployed at", n.DeployedAt.Format("2006-01-02 15:04:05"))
	}

	return c.sendMessage(ctx, tkn, channelID, msgOptions.Build())
}

func (c *Client) sendMessage(ctx context.Context, tkn model.SlackToken, channelID string, msgOptions []slack.MsgOption) error {
	client := slack.New(tkn.String())

	if _, _, err := client.PostMessageContext(ctx, channelID, msgOptions...); err != nil {
		switch err.Error() {
		case errInvalidAuth:
			return svcerrors.NewSlackClientUnauthorizedError().Wrap(err)
		case errChannelNotFound:
			return svcerrors.NewSlackChannelNotFoundError().Wrap(err).
				WithMessage(fmt.Sprintf("slack channel (ID: %s) not found", channelID))
		default:
			return fmt.Errorf("sending message to Slack channel (ID: %s): %w", channelID, err)
		}
	}

	return nil
}

package slack

import (
	"context"
	"fmt"

	svcerrors "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/slack-go/slack"
	"go.strv.io/background"
	"go.strv.io/background/task"
)

const (
	errInvalidAuth     = "invalid_auth"
	errChannelNotFound = "channel_not_found"
)

type Client struct {
	taskManager *background.Manager
}

func NewClient(manager *background.Manager) *Client {
	return &Client{
		taskManager: manager,
	}
}

func (c *Client) SendReleaseNotificationAsync(ctx context.Context, token, channelID string, n model.ReleaseNotification) {
	c.sendMessageAsync(ctx, token, channelID, c.buildReleaseNotificationMsgOptions(n))
}

func (c *Client) SendReleaseNotification(ctx context.Context, token, channelID string, n model.ReleaseNotification) error {
	return c.sendMessage(ctx, token, channelID, c.buildReleaseNotificationMsgOptions(n))
}

func (c *Client) buildReleaseNotificationMsgOptions(n model.ReleaseNotification) []slack.MsgOption {
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

	return msgOptions.Build()
}

func (c *Client) sendMessageAsync(ctx context.Context, token, channelID string, msgOptions []slack.MsgOption) {
	t := task.Task{
		Type: task.TypeOneOff,
		Meta: task.Metadata{
			"task": "sending slack message",
		},
		Fn: func(ctx context.Context) error {
			return c.sendMessage(ctx, token, channelID, msgOptions)
		},
	}

	c.taskManager.RunTask(ctx, t)
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

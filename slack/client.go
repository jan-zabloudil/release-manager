package slack

import (
	"context"
	"fmt"

	"release-manager/service/model"

	"github.com/slack-go/slack"
	"go.strv.io/background"
	"go.strv.io/background/task"
)

type Client struct {
	taskManager *background.Manager
}

func NewClient(manager *background.Manager) *Client {
	return &Client{
		taskManager: manager,
	}
}

func (c *Client) SendReleaseNotification(ctx context.Context, token, channelID string, n model.ReleaseNotification) {
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

	c.sendMessageAsync(ctx, token, channelID, msgOptions.Build())
}

func (c *Client) sendMessageAsync(ctx context.Context, token, channelID string, msgOptions []slack.MsgOption) {
	client := slack.New(token)

	t := task.Task{
		Type: task.TypeOneOff,
		Meta: task.Metadata{
			"task": "sending slack message",
		},
		Fn: func(ctx context.Context) error {
			if _, _, err := client.PostMessageContext(ctx, channelID, msgOptions...); err != nil {
				return fmt.Errorf("failed to send message to Slack channel (ID %s): %w", channelID, err)
			}

			return nil
		},
	}

	c.taskManager.RunTask(ctx, t)
}

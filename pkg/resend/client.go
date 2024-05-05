package resend

import (
	"context"

	"release-manager/config"
	"release-manager/pkg/emailerrors"

	"github.com/resend/resend-go/v2"
	"go.strv.io/background"
	"go.strv.io/background/task"
)

type Client struct {
	taskManager *background.Manager
	client      *resend.Client
	reqBuilder  *EmailRequestBuilder
}

func NewClient(manager *background.Manager, cfg config.ResendConfig) *Client {
	return &Client{
		taskManager: manager,
		client:      resend.NewClient(cfg.APIKey),
		reqBuilder:  NewEmailRequestBuilder(cfg),
	}
}

func (c *Client) SendEmailAsync(ctx context.Context, subject, text, html string, recipients ...string) {
	req := c.reqBuilder.
		SetRecipients(recipients).
		SetSubject(subject).
		SetText(text).
		SetHTML(html).
		Build()

	t := task.Task{
		Type: task.TypeOneOff,
		Meta: task.Metadata{
			"task": "sending email",
		},
		Fn: func(ctx context.Context) error {
			_, err := c.client.Emails.SendWithContext(ctx, req)
			if err != nil {
				return emailerrors.NewEmailSendingFailedError().Wrap(err)
			}

			return nil
		},
	}

	c.taskManager.RunTask(ctx, t)
}

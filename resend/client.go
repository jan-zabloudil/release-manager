package resend

import (
	"context"
	"fmt"

	"release-manager/config"
	"release-manager/service/model"

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

func (c *Client) SendEmailAsync(ctx context.Context, email model.Email) {
	req := c.reqBuilder.
		SetRecipients(email.Recipients).
		SetSubject(email.Subject).
		SetText(email.Text).
		SetHTML(email.HTML).
		Build()

	t := task.Task{
		Type: task.TypeOneOff,
		Meta: task.Metadata{
			"task": "sending email",
		},
		Fn: func(ctx context.Context) error {
			_, err := c.client.Emails.SendWithContext(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to send email via Resend: %w", err)
			}

			return nil
		},
	}

	c.taskManager.RunTask(ctx, t)
}

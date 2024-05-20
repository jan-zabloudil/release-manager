package resend

import (
	"context"
	"fmt"
	"log/slog"

	"release-manager/config"
	"release-manager/service/model"

	"github.com/resend/resend-go/v2"
	"go.strv.io/background"
	"go.strv.io/background/task"
)

type Client struct {
	taskManager  *background.Manager
	client       *resend.Client
	clientSvcCfg config.ClientServiceConfig
	reqBuilder   *EmailRequestBuilder
}

func NewClient(manager *background.Manager, resendCfg config.ResendConfig, clientSvcCfg config.ClientServiceConfig) *Client {
	return &Client{
		taskManager:  manager,
		client:       resend.NewClient(resendCfg.APIKey),
		clientSvcCfg: clientSvcCfg,
		reqBuilder:   NewEmailRequestBuilder(resendCfg),
	}
}

func (c *Client) SendProjectInvitationEmailAsync(
	ctx context.Context,
	data model.ProjectInvitationEmailData,
	recipient string,
) {
	parsedTmpl, err := ParseProjectInvitationTemplate(data, c.clientSvcCfg)
	if err != nil {
		slog.Error("failed to parse project invitation template", "error", err)
		return
	}

	c.sendEmailAsync(ctx, parsedTmpl.Subject, parsedTmpl.Text, parsedTmpl.HTML, recipient)
}

func (c *Client) sendEmailAsync(ctx context.Context, subject, text, html string, recipients ...string) {
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
				return fmt.Errorf("failed to send email via Resend: %w", err)
			}

			return nil
		},
	}

	c.taskManager.RunTask(ctx, t)
}

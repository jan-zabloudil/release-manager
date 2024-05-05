package resend

import (
	"release-manager/config"

	"github.com/resend/resend-go/v2"
)

type EmailRequestBuilder struct {
	request              *resend.SendEmailRequest
	sendToRealRecipients bool
}

func NewEmailRequestBuilder(cfg config.ResendConfig) *EmailRequestBuilder {
	return &EmailRequestBuilder{
		request: &resend.SendEmailRequest{
			From: cfg.Sender,
			To:   []string{cfg.TestRecipient},
		},
		sendToRealRecipients: cfg.SendToRealRecipients,
	}
}

func (b *EmailRequestBuilder) SetSubject(subject string) *EmailRequestBuilder {
	b.request.Subject = subject
	return b
}

func (b *EmailRequestBuilder) SetText(text string) *EmailRequestBuilder {
	b.request.Text = text
	return b
}

func (b *EmailRequestBuilder) SetHTML(html string) *EmailRequestBuilder {
	b.request.Html = html
	return b
}

func (b *EmailRequestBuilder) SetRecipients(r []string) *EmailRequestBuilder {
	if b.sendToRealRecipients {
		b.request.To = r
	}
	return b
}

func (b *EmailRequestBuilder) Build() *resend.SendEmailRequest {
	return b.request
}

package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"

	"github.com/resend/resend-go/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Config struct {
	TestingMode   bool
	ApiKey        string
	Sender        string
	TestRecipient string
}

type Mailer struct {
	resend        *resend.Client
	testingMode   bool
	sender        string
	testRecipient string
}

func New(cfg Config) *Mailer {
	return &Mailer{
		resend:        resend.NewClient(cfg.ApiKey),
		testingMode:   cfg.TestingMode,
		sender:        cfg.Sender,
		testRecipient: cfg.TestRecipient,
	}
}

func (m *Mailer) buildEmailRequest(recipients []string, templateFile string) (*resend.SendEmailRequest, error) {
	files := []string{
		"templates/base.html",
		"templates/partials/*.html",
		"templates/" + templateFile,
	}

	tmpl, err := template.New("email").ParseFS(templateFS, files...)

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", nil)
	if err != nil {
		return nil, err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", nil)
	if err != nil {
		return nil, err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", nil)
	if err != nil {
		return nil, err
	}

	if m.testingMode {
		recipients = []string{m.testRecipient}
	}

	req := &resend.SendEmailRequest{
		To:      recipients,
		From:    m.sender,
		Text:    plainBody.String(),
		Html:    htmlBody.String(),
		Subject: fmt.Sprintf("%s | ReleaseManager", subject.String()),
	}

	return req, nil
}

func (m *Mailer) sendEmail(ctx context.Context, req *resend.SendEmailRequest) error {
	sent, err := m.resend.Emails.SendWithContext(ctx, req)
	if err != nil {
		return err
	}

	slog.Debug("email sent", "email id", sent.Id)
	return nil
}

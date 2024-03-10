package mailer

import "context"

func (m *Mailer) SendTestEmail(ctx context.Context, recipients []string) error {
	params, err := m.buildEmailRequest(recipients, "test.html")
	if err != nil {
		return err
	}

	if err := m.sendEmail(ctx, params); err != nil {
		return err
	}

	return nil
}

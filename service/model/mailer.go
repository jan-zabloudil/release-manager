package model

import "context"

type Mailer interface {
	SendTestEmail(ctx context.Context, recipients []string) error
}

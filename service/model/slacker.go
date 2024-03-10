package model

import "context"

type Slacker interface {
	PostTestMessage(ctx context.Context, channelID string) error
}

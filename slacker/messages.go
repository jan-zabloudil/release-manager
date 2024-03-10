package slacker

import (
	"context"
	"log/slog"
)

func (s *Slacker) PostTestMessage(ctx context.Context, channelID string) error {
	text := ":rocket: *Post message test* :rocket:\n\nHey, it is working!"

	return s.sendTextMessage(ctx, channelID, text)
}

func (s *SilentSlacker) PostTestMessage(_ context.Context, channelID string) error {
	slog.Debug("no message send to slack because silent slacker is used", "channelId", channelID)

	return nil
}

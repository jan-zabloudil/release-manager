package slack

import (
	"context"

	svcmodel "release-manager/service/model"
	"release-manager/slack/model"
)

func (s *Slack) PostReleaseMessage(ctx context.Context, p svcmodel.Project, app svcmodel.App, rls svcmodel.Release) error {
	// Project will hold information about release message
	// Project will also hold information about channelID
	// MessageTemplate could provide a method that would return map of :
	// If information should be included

	title := ":rocket: *New release* :rocket:"
	text := "Hi everyone,\\nI’m thrilled to announce new release!"

	msg := model.NewMessage(title, text)

	msg.WithField("Project", p.Name)

	return nil
}

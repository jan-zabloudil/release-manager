package slack

import (
	"context"
	"fmt"

	svcmodel "release-manager/service/model"
	"release-manager/slack/model"
)

func (s *Slack) PostReleaseMessage(ctx context.Context, p svcmodel.Project, app svcmodel.App, rls svcmodel.Release) error {
	rm := svcmodel.ReleaseMessage{
		Title: ":rocket: *New release* :rocket:",
		Text:  "Hi everyone,\\nI’m thrilled to announce new release!",
		Includes: svcmodel.Includes{
			ProjectName: true,
			AppName:     true,
			ReleaseName: true,
			Changelog:   true,
		},
	}

	msg := model.NewMessage(rm.Title, rm.Text)

	if rm.Includes.ProjectName {
		msg.WithField("Project", p.Name)
	}
	if rm.Includes.AppName {
		msg.WithField("AppName", app.Name)
	}
	if rm.Includes.ReleaseName {
		msg.WithField("Release", rls.Title)
	}
	if rm.Includes.Changelog {
		msg.WithField("Changelog", rls.ChangeLog)
	}
	if rm.Includes.Deployments {
		var dplMsg string

		// todo ošetřit, když není environment vyplněn
		if rls.Deployments.Dev {
			dplMsg += fmt.Sprintf("Dev: %s", app.Environments.DevURL)
		}
		if rls.Deployments.Dev {
			dplMsg += fmt.Sprintf("Staging: %s", app.Environments.StgURL)
		}
		if rls.Deployments.Dev {
			dplMsg += fmt.Sprintf("Staging: %s", app.Environments.PrdURL)
		}
	}

	if err := s.sendMessageWithAttachments(ctx, "C065E66TZ36", msg); err != nil {
		return err
	}

	return nil
}

package slack

import (
	"context"
	"fmt"

	svcmodel "release-manager/service/model"
	"release-manager/slack/model"
)

func (s *Slack) PostReleaseMessage(ctx context.Context, p svcmodel.Project, app svcmodel.App, rls svcmodel.Release) error {
	rm := svcmodel.ReleaseMessage{
		Title: ":rocket: New release :rocket:",
		Text:  "Hi everyone,\nI’m thrilled to announce new release!",
		Includes: svcmodel.Includes{
			ProjectName: true,
			AppName:     true,
			ReleaseName: true,
			Changelog:   true,
			Deployments: true,
		},
	}

	msg := model.NewReleaseMessage(rm.Title, rm.Text)

	if rm.Includes.ProjectName {
		msg.WithField("Project", p.Name)
	}
	if rm.Includes.AppName {
		msg.WithField("App", app.Name)
	}
	if rm.Includes.ReleaseName {
		msg.WithField("Release", rls.Title)
	}
	if rm.Includes.Deployments {
		var dplMsg string

		fmt.Println("beforeee")
		fmt.Println(app.Environments.DevURL)

		if rls.Deployments.Dev {
			if app.Environments.DevURL != nil {
				dplMsg += fmt.Sprintf("Dev: %s\n", app.Environments.DevURL.String())
			} else {
				dplMsg += "Dev\n"
			}
		}
		if rls.Deployments.Stg {
			if app.Environments.StgURL != nil {
				dplMsg += fmt.Sprintf("Staging: %s\n", app.Environments.StgURL.String())
			} else {
				dplMsg += "Staging\n"
			}
		}
		if rls.Deployments.Prd {
			if app.Environments.PrdURL != nil {
				dplMsg += fmt.Sprintf("Prod: %s\n", app.Environments.PrdURL.String())
			} else {
				dplMsg += "Prod\n"
			}
		}

		if dplMsg != "" {
			msg.WithField("Deployments", dplMsg)
		}
	}
	if rm.Includes.Changelog && rls.ChangeLog != "" {
		msg.WithField("Changelog", rls.ChangeLog)
	}

	if err := s.sendMessage(ctx, p.Notifications.SlackChannelID, msg.Options()); err != nil {
		return err
	}

	return nil
}

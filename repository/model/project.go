package model

import (
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/relvacode/iso8601"
)

type ProjectNotifications struct {
	SlackChannelID         string `json:"slack_channel_id"`
	ReleaseMessageTemplate string `json:"release_message_template"`
}

type ProjectInput struct {
	Name                   string `json:"name"`
	Description            string `json:"description"`
	SlackChannelID         string `json:"slack_channel_id"`
	ReleaseMessageTemplate string `json:"release_message_template"`
}

type ProjectResponse struct {
	ID                     uuid.UUID    `json:"id"`
	Name                   string       `json:"name"`
	Description            string       `json:"description"`
	SlackChannelID         string       `json:"slack_channel_id"`
	ReleaseMessageTemplate string       `json:"release_message_template"`
	CreatedAt              iso8601.Time `json:"created_at"`
	UpdatedAt              iso8601.Time `json:"updated_at"`
}

func ToProjectDBInput(p svcmodel.Project) ProjectInput {
	return ProjectInput{
		Name:                   p.Name,
		Description:            p.Description,
		SlackChannelID:         p.Notifications.SlackChannelID,
		ReleaseMessageTemplate: p.Notifications.ReleaseMessageTemplate,
	}
}

func ToSvcProject(r ProjectResponse) svcmodel.Project {
	return svcmodel.Project{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Notifications: svcmodel.ProjectNotifications{
			SlackChannelID:         r.SlackChannelID,
			ReleaseMessageTemplate: r.ReleaseMessageTemplate,
		},
		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}
}

func ToSvcProjects(pr []ProjectResponse) []svcmodel.Project {
	svcProjects := make([]svcmodel.Project, 0, len(pr))
	for _, p := range pr {
		svcProjects = append(svcProjects, ToSvcProject(p))
	}

	return svcProjects
}

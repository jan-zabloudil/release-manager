package model

import (
	"context"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectService interface {
	Create(ctx context.Context, p svcmodel.Project) (svcmodel.Project, error)
	Get(ctx context.Context, id uuid.UUID) (svcmodel.Project, error)
	ListAll(ctx context.Context) ([]svcmodel.Project, error)
	Update(ctx context.Context, project svcmodel.Project) (svcmodel.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectNotifications struct {
	SlackChannelID         string `json:"slack_channel_id"`
	ReleaseMessageTemplate string `json:"release_message_template"`
}

type Project struct {
	ID            uuid.UUID            `json:"id"`
	Name          string               `json:"name" validate:"required"`
	Description   string               `json:"description"`
	Notifications ProjectNotifications `json:"notifications"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

func ToSvcProject(p Project) (sv svcmodel.Project) {
	return svcmodel.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Notifications: svcmodel.ProjectNotifications{
			SlackChannelID:         p.Notifications.SlackChannelID,
			ReleaseMessageTemplate: p.Notifications.ReleaseMessageTemplate,
		},
	}
}

func ToNetProject(p svcmodel.Project) Project {
	return Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Notifications: ProjectNotifications{
			SlackChannelID:         p.Notifications.SlackChannelID,
			ReleaseMessageTemplate: p.Notifications.ReleaseMessageTemplate,
		},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToNetProjects(svcProjects []svcmodel.Project) []Project {
	netProjects := make([]Project, 0, len(svcProjects))
	for _, svcProject := range svcProjects {
		netProjects = append(netProjects, ToNetProject(svcProject))
	}

	return netProjects
}

type ProjectNotificationsPatch struct {
	SlackChannelID         *string `json:"slack_channel_id"`
	ReleaseMessageTemplate *string `json:"release_message_template"`
}

type ProjectPatch struct {
	Name          *string                    `json:"name"`
	Description   *string                    `json:"description"`
	Notifications *ProjectNotificationsPatch `json:"notifications"`
}

func PatchToNetProject(input ProjectPatch, p Project) Project {
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Description != nil {
		p.Description = *input.Description
	}
	if input.Notifications != nil {
		if input.Notifications.SlackChannelID != nil {
			p.Notifications.SlackChannelID = *input.Notifications.SlackChannelID
		}
		if input.Notifications.ReleaseMessageTemplate != nil {
			p.Notifications.ReleaseMessageTemplate = *input.Notifications.ReleaseMessageTemplate
		}
	}

	return p
}

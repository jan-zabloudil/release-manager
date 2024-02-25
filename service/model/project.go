package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Insert(ctx context.Context, p Project) (Project, error)
	Read(ctx context.Context, id uuid.UUID) (Project, error)
	ReadAll(ctx context.Context) ([]Project, error)
	Update(ctx context.Context, p Project) (Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectNotifications struct {
	SlackChannelID         string
	ReleaseMessageTemplate string
}

// TODO Add ProjectNotifications to validate if Slack message can be sent

type Project struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Notifications ProjectNotifications
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

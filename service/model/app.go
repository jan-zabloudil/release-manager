package model

import (
	"context"
	"net/url"
	"time"

	svcerr "release-manager/service/errors"

	"github.com/google/uuid"
)

type AppRepository interface {
	Insert(ctx context.Context, app App) (App, error)
	Read(ctx context.Context, id uuid.UUID) (App, error)
	ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]App, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, app App) (App, error)
}

type App struct {
	ID           uuid.UUID
	ProjectID    uuid.UUID
	Name         string
	Description  string
	Environments Environments
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Environments struct {
	DevURL EnvURL
	StgURL EnvURL
	PrdURL EnvURL
}

type envURL struct {
	url *url.URL
}

type EnvURL interface {
	String() string
	IsEnvURL() bool
}

func (u *envURL) String() string {
	return u.url.String()
}

// TODO how to better create unique interface?
func (u *envURL) IsEnvURL() bool {
	return true
}

func NewEnvURL(raw string) (EnvURL, error) {
	parsedURL, err := url.ParseRequestURI(raw)
	if err != nil {
		return nil, svcerr.ErrAppEnvURLInvalid
	}

	if !parsedURL.IsAbs() {
		return nil, svcerr.ErrAppEnvURLInvalid
	}

	return &envURL{url: parsedURL}, nil
}

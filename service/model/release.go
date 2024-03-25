package model

import (
	"context"
	"time"

	svcerr "release-manager/service/errors"

	"github.com/google/uuid"
)

type ReleaseRepository interface {
	Insert(ctx context.Context, r Release) (Release, error)
	ReadAllForApp(ctx context.Context, appID uuid.UUID) ([]Release, error)
	Read(ctx context.Context, id uuid.UUID) (Release, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, r Release) (Release, error)
}

type Slack interface {
	PostReleaseMessage(ctx context.Context, p Project, app App, rls Release) error
}

type AppService interface {
	Get(ctx context.Context, id uuid.UUID) (App, error)
}

type ProjectService interface {
	Get(ctx context.Context, id uuid.UUID) (Project, error)
}

type Release struct {
	ID              uuid.UUID
	AppID           uuid.UUID
	SourceCode      SourceCode
	Deployments     Deployments
	Title           string
	ChangeLog       string
	CreatedByUserID uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type sourceCode struct {
	tag             string
	targetCommitIsh string
}

type SourceCode interface {
	Tag() string
	TargetCommitIsh() string
}

func NewSourceCode(tag, targetCommitIsh string) (SourceCode, error) {
	if tag == "" {
		return nil, svcerr.ErrInvalidTag
	}

	return &sourceCode{
		tag:             tag,
		targetCommitIsh: targetCommitIsh,
	}, nil
}

func (s *sourceCode) Tag() string {
	return s.tag
}

func (s *sourceCode) TargetCommitIsh() string {
	return s.targetCommitIsh
}

type Deployments struct {
	Dev bool
	Stg bool
	Prd bool
}

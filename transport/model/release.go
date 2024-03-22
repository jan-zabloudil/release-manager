package model

import (
	"context"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ReleaseService interface {
	Create(ctx context.Context, r svcmodel.Release) (svcmodel.Release, error)
	GetAllForApp(ctx context.Context, appID uuid.UUID) ([]svcmodel.Release, error)
	Get(ctx context.Context, id uuid.UUID) (svcmodel.Release, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, rls svcmodel.Release) (svcmodel.Release, error)
}

type Release struct {
	ID              uuid.UUID   `json:"id"`
	SourceCode      SourceCode  `json:"source_code" validate:"source_code_required"`
	Deployments     Deployments `json:"deployments"`
	Title           *string     `json:"title"`
	ChangeLog       *string     `json:"changelog"`
	CreatedByUserID uuid.UUID   `json:"created_by_user_id"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type ReleasePatch struct {
	SourceCode  SourceCode  `json:"source_code" validate:"source_code_if_present"`
	Deployments Deployments `json:"deployments"`
	Title       *string     `json:"title"`
	ChangeLog   *string     `json:"changelog"`
}

type SourceCode struct {
	Tag             *string `json:"tag"`
	TargetCommitIsh *string `json:"target_commitish"`
}

type Deployments struct {
	Dev *bool `json:"dev"`
	Stg *bool `json:"staging"`
	Prd *bool `json:"production"`
}

func NewSvcRelease(appID uuid.UUID, sc SourceCode, title, changeLog *string, dev, stg, prd *bool, createdByUserID uuid.UUID) (svcmodel.Release, error) {
	var r svcmodel.Release

	r.ID = uuid.New()
	r.AppID = appID
	r.CreatedAt = time.Now()
	r.CreatedByUserID = createdByUserID

	return ToSvcRelease(r, sc, title, changeLog, dev, stg, prd)
}

func ToSvcRelease(r svcmodel.Release, sc SourceCode, title, changeLog *string, dev, stg, prd *bool) (svcmodel.Release, error) {
	if sc.Tag != nil && sc.TargetCommitIsh != nil {
		sourceCode, err := svcmodel.NewSourceCode(*sc.Tag, *sc.TargetCommitIsh)
		if err != nil {
			return svcmodel.Release{}, err
		}
		r.SourceCode = sourceCode
	}

	if dev != nil {
		r.Deployments.Dev = *dev
	}
	if stg != nil {
		r.Deployments.Stg = *stg
	}
	if prd != nil {
		r.Deployments.Prd = *prd
	}
	if title != nil {
		r.Title = *title
	}
	if changeLog != nil {
		r.ChangeLog = *changeLog
	}
	r.UpdatedAt = time.Now()

	return r, nil
}

func ToNetRelease(
	id uuid.UUID,
	tag string,
	targetCommitIsh string,
	dev bool,
	stg bool,
	prd bool,
	title string,
	changelog string,
	createdByUserID uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
) Release {
	return Release{
		ID: id,
		SourceCode: SourceCode{
			Tag:             &tag,
			TargetCommitIsh: &targetCommitIsh,
		},
		Deployments: Deployments{
			Dev: &dev,
			Stg: &stg,
			Prd: &prd,
		},
		Title:           &title,
		ChangeLog:       &changelog,
		CreatedByUserID: createdByUserID,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

func ToNetReleases(releases []svcmodel.Release) []Release {
	r := make([]Release, 0, len(releases))
	for _, svcApp := range releases {
		r = append(r, ToNetRelease(
			svcApp.ID,
			svcApp.SourceCode.Tag(),
			svcApp.SourceCode.TargetCommitIsh(),
			svcApp.Deployments.Dev,
			svcApp.Deployments.Stg,
			svcApp.Deployments.Prd,
			svcApp.Title,
			svcApp.ChangeLog,
			svcApp.CreatedByUserID,
			svcApp.CreatedAt,
			svcApp.UpdatedAt,
		))
	}

	return r
}

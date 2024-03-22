package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type Release struct {
	ID              uuid.UUID   `json:"id"`
	AppID           uuid.UUID   `json:"app_id"`
	SourceCode      SourceCode  `json:"source_code"`
	Deployments     Deployments `json:"deployments"`
	Title           string      `json:"title"`
	ChangeLog       string      `json:"changelog"`
	CreatedByUserID uuid.UUID   `json:"created_by_user_id"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type SourceCode struct {
	Tag             string `json:"tag"`
	TargetCommitIsh string `json:"target_commitish"`
}

type Deployments struct {
	Dev bool `json:"dev"`
	Stg bool `json:"stg"`
	Prd bool `json:"prd"`
}

func ToDBRelease(
	id uuid.UUID,
	appID uuid.UUID,
	createdByUserID uuid.UUID,
	tag string,
	targetCommitIsh string,
	dev bool,
	stg bool,
	prd bool,
	title string,
	changelog string,
	createdAt time.Time,
	updatedAt time.Time,
) Release {
	return Release{
		ID:    id,
		AppID: appID,
		SourceCode: SourceCode{
			Tag:             tag,
			TargetCommitIsh: targetCommitIsh,
		},
		Deployments: Deployments{
			Dev: dev,
			Stg: stg,
			Prd: prd,
		},
		Title:           title,
		ChangeLog:       changelog,
		CreatedByUserID: createdByUserID,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

func ToSvcRelease(
	id uuid.UUID,
	appID uuid.UUID,
	createdByUserID uuid.UUID,
	tag string,
	targetCommitIsh string,
	dev bool,
	stg bool,
	prd bool,
	title string,
	changelog string,
	createdAt time.Time,
	updatedAt time.Time,
) (svcmodel.Release, error) {
	sourceCode, err := svcmodel.NewSourceCode(tag, targetCommitIsh)
	if err != nil {
		return svcmodel.Release{}, err
	}

	return svcmodel.Release{
		ID:         id,
		AppID:      appID,
		SourceCode: sourceCode,
		Deployments: svcmodel.Deployments{
			Dev: dev,
			Stg: stg,
			Prd: prd,
		},
		Title:           title,
		ChangeLog:       changelog,
		CreatedByUserID: createdByUserID,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

func ToSvcReleases(releases []Release) ([]svcmodel.Release, error) {
	r := make([]svcmodel.Release, 0, len(releases))
	for _, release := range releases {
		svcRelease, err := ToSvcRelease(
			release.ID,
			release.AppID,
			release.CreatedByUserID,
			release.SourceCode.Tag,
			release.SourceCode.TargetCommitIsh,
			release.Deployments.Dev,
			release.Deployments.Stg,
			release.Deployments.Prd,
			release.Title,
			release.ChangeLog,
			release.CreatedAt,
			release.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		r = append(r, svcRelease)
	}

	return r, nil
}

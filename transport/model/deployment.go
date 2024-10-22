package model

import (
	"fmt"
	"strconv"
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
)

type CreateDeploymentInput struct {
	ReleaseID     id.Release     `json:"release_id"`
	EnvironmentID id.Environment `json:"environment_id"`
}

type Deployment struct {
	ID                    id.Deployment  `json:"id"`
	ReleaseID             id.Release     `json:"release_id"`
	ReleaseTitle          string         `json:"release_title"`
	EnvironmentID         id.Environment `json:"environment_id"`
	EnvironmentName       string         `json:"environment_name"`
	EnvironmentServiceURL string         `json:"environment_service_url"`
	DeployedByUserID      id.AuthUser    `json:"deployed_by_user_id"`
	DeployedAt            time.Time      `json:"deployed_at"`
}

func ToSvcCreateDeploymentInput(input CreateDeploymentInput) svcmodel.CreateDeploymentInput {
	return svcmodel.CreateDeploymentInput{
		ReleaseID:     input.ReleaseID,
		EnvironmentID: input.EnvironmentID,
	}
}

func ToSvcDeploymentFilterParams(releaseIDParam, environmentIDParam, latestOnlyParam string) (svcmodel.DeploymentFilterParams, error) {
	var releaseID *id.Release
	var environmentID *id.Environment
	var latestOnly *bool

	if releaseIDParam != "" {
		var parsedID id.Release
		if err := parsedID.FromString(releaseIDParam); err != nil {
			return svcmodel.DeploymentFilterParams{}, fmt.Errorf("invalid uuid provided for release id: %w", err)
		}
		releaseID = &parsedID
	}

	if environmentIDParam != "" {
		var parsedID id.Environment
		if err := parsedID.FromString(environmentIDParam); err != nil {
			return svcmodel.DeploymentFilterParams{}, fmt.Errorf("invalid uuid provided for environment id: %w", err)
		}
		environmentID = &parsedID
	}

	if latestOnlyParam != "" {
		latestOnlyValue, err := strconv.ParseBool(latestOnlyParam)
		if err != nil {
			return svcmodel.DeploymentFilterParams{}, fmt.Errorf("invalid boolean provided for latest only: %w", err)
		}
		latestOnly = &latestOnlyValue
	}

	return svcmodel.DeploymentFilterParams{
		ReleaseID:     releaseID,
		EnvironmentID: environmentID,
		LatestOnly:    latestOnly,
	}, nil
}

func ToDeployment(dpl svcmodel.Deployment) Deployment {
	return Deployment{
		ID:                    dpl.ID,
		ReleaseID:             dpl.Release.ID,
		ReleaseTitle:          dpl.Release.ReleaseTitle,
		EnvironmentID:         dpl.Environment.ID,
		EnvironmentName:       dpl.Environment.Name,
		EnvironmentServiceURL: dpl.Environment.ServiceURL.String(),
		DeployedByUserID:      dpl.DeployedByUserID,
		DeployedAt:            dpl.DeployedAt,
	}
}

func ToDeployments(dpls []svcmodel.Deployment) []Deployment {
	d := make([]Deployment, 0, len(dpls))
	for _, dpl := range dpls {
		d = append(d, ToDeployment(dpl))
	}
	return d
}

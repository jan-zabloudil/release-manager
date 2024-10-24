package model

import (
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

type ListDeploymentsParams struct {
	ProjectID     id.Project      `param:"path=project_id"`
	ReleaseID     *id.Release     `param:"query=release_id"`
	EnvironmentID *id.Environment `param:"query=environment_id"`
	LatestOnly    *bool           `param:"query=latest_only"`
}

func ToSvcCreateDeploymentInput(input CreateDeploymentInput) svcmodel.CreateDeploymentInput {
	return svcmodel.CreateDeploymentInput{
		ReleaseID:     input.ReleaseID,
		EnvironmentID: input.EnvironmentID,
	}
}

func ToSvcListDeploymentsFilterParams(p ListDeploymentsParams) svcmodel.ListDeploymentsFilterParams {
	return svcmodel.ListDeploymentsFilterParams{
		ReleaseID:     p.ReleaseID,
		EnvironmentID: p.EnvironmentID,
		LatestOnly:    p.LatestOnly,
	}
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

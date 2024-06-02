package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type CreateDeploymentInput struct {
	ReleaseID     uuid.UUID `json:"release_id"`
	EnvironmentID uuid.UUID `json:"environment_id"`
}

type Deployment struct {
	ID               uuid.UUID `json:"id"`
	ReleaseID        uuid.UUID `json:"release_id"`
	ReleaseTitle     string    `json:"release_title"`
	EnvironmentID    uuid.UUID `json:"environment_id"`
	EnvironmentName  string    `json:"environment_name"`
	DeployedByUserID uuid.UUID `json:"deployed_by_user_id"`
	DeployedAt       time.Time `json:"deployed_at"`
}

func ToSvcCreateDeploymentInput(input CreateDeploymentInput) svcmodel.CreateDeploymentInput {
	return svcmodel.CreateDeploymentInput{
		ReleaseID:     input.ReleaseID,
		EnvironmentID: input.EnvironmentID,
	}
}

func ToDeployment(dpl svcmodel.Deployment) Deployment {
	return Deployment{
		ID:               dpl.ID,
		ReleaseID:        dpl.Release.ID,
		ReleaseTitle:     dpl.Release.ReleaseTitle,
		EnvironmentID:    dpl.Environment.ID,
		EnvironmentName:  dpl.Environment.Name,
		DeployedByUserID: dpl.DeployedByUserID,
		DeployedAt:       dpl.DeployedAt,
	}
}

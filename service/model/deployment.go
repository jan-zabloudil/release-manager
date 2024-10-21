package model

import (
	"errors"
	"time"

	"release-manager/pkg/id"

	"github.com/google/uuid"
)

var (
	errReleaseIDRequired     = errors.New("release id is required")
	errEnvironmentIDRequired = errors.New("environment id is required")
)

type CreateDeploymentInput struct {
	ReleaseID     uuid.UUID
	EnvironmentID uuid.UUID
}

func (i CreateDeploymentInput) Validate() error {
	if i.ReleaseID == uuid.Nil {
		return errReleaseIDRequired
	}
	if i.EnvironmentID == uuid.Nil {
		return errEnvironmentIDRequired
	}
	return nil
}

type Deployment struct {
	ID               id.Deployment
	Release          Release
	Environment      Environment
	DeployedByUserID id.AuthUser
	DeployedAt       time.Time
}

func NewDeployment(rls Release, env Environment, deployedByUserID id.AuthUser) Deployment {
	return Deployment{
		ID:               id.NewDeployment(),
		Release:          rls,
		Environment:      env,
		DeployedByUserID: deployedByUserID,
		DeployedAt:       time.Now(),
	}
}

type DeploymentFilterParams struct {
	ReleaseID     *uuid.UUID
	EnvironmentID *uuid.UUID
	LatestOnly    *bool
}

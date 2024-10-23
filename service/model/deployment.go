package model

import (
	"errors"
	"time"

	"release-manager/pkg/id"
)

var (
	errReleaseIDRequired     = errors.New("release id is required")
	errEnvironmentIDRequired = errors.New("environment id is required")
)

type CreateDeploymentInput struct {
	ReleaseID     id.Release
	EnvironmentID id.Environment
}

func (i CreateDeploymentInput) Validate() error {
	if i.ReleaseID.IsNil() {
		return errReleaseIDRequired
	}
	if i.EnvironmentID.IsNil() {
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

type ListDeploymentsFilterParams struct {
	ReleaseID     *id.Release
	EnvironmentID *id.Environment
	LatestOnly    *bool
}

package model

import (
	"errors"
	"time"

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
	ID               uuid.UUID
	Release          Release
	Environment      Environment
	DeployedByUserID uuid.UUID
	DeployedAt       time.Time
}

func NewDeployment(rls Release, env Environment, deployedByUserID uuid.UUID) Deployment {
	return Deployment{
		ID:               uuid.New(),
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

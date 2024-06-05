package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeploymentInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateDeploymentInput
		wantErr bool
	}{
		{
			name: "Valid Input",
			input: CreateDeploymentInput{
				ReleaseID:     uuid.New(),
				EnvironmentID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "Invalid Input - No ReleaseID",
			input: CreateDeploymentInput{
				EnvironmentID: uuid.New(),
			},
			wantErr: true,
		},
		{
			name: "Invalid Input - No EnvironmentID",
			input: CreateDeploymentInput{
				ReleaseID: uuid.New(),
			},
			wantErr: true,
		},
		{
			name:    "Invalid Input - No ReleaseID and EnvironmentID",
			input:   CreateDeploymentInput{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

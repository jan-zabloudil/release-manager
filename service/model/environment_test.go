package model

import (
	"testing"

	"release-manager/pkg/pointer"
	"release-manager/pkg/urlx"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_NewEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		creation CreateEnvironmentInput
		wantErr  bool
	}{
		{
			name: "Valid Environment",
			creation: CreateEnvironmentInput{
				ProjectID:     uuid.New(),
				Name:          "dev",
				ServiceRawURL: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Invalid Environment - not absolute service url",
			creation: CreateEnvironmentInput{
				ProjectID:     uuid.New(),
				Name:          "dev",
				ServiceRawURL: "example.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Environment - empty name",
			creation: CreateEnvironmentInput{
				ProjectID:     uuid.New(),
				Name:          "",
				ServiceRawURL: "http://example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEnvironment(tt.creation)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnvironment_Validate(t *testing.T) {
	tests := []struct {
		name    string
		env     Environment
		wantErr bool
	}{
		{
			name: "Valid Environment",
			env: Environment{
				Name: "Test Environment",
			},
			wantErr: false,
		},
		{
			name:    "Invalid Environment - Empty Name",
			env:     Environment{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.env.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnvironment_Update(t *testing.T) {
	tests := []struct {
		name    string
		env     Environment
		update  UpdateEnvironmentInput
		wantErr bool
	}{
		{
			name: "Valid Update",
			env: Environment{
				Name: "Old Name",
			},
			update: UpdateEnvironmentInput{
				Name:          pointer.StringPtr("New name"),
				ServiceRawURL: pointer.StringPtr("https://new.example.com"),
			},
			wantErr: false,
		},
		{
			name: "Invalid Update - not absolute service url",
			env: Environment{
				Name: "Old Name",
			},
			update: UpdateEnvironmentInput{
				Name:          pointer.StringPtr("New name"),
				ServiceRawURL: pointer.StringPtr("relative-url"),
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - empty name",
			env: Environment{
				Name: "Old Name",
			},
			update: UpdateEnvironmentInput{
				Name:          pointer.StringPtr(""),
				ServiceRawURL: pointer.StringPtr("https://example.com"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.env.Update(tt.update)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, *tt.update.Name, tt.env.Name)
				assert.Equal(t, *tt.update.ServiceRawURL, tt.env.ServiceURL.String())
			}
		})
	}
}

func TestIsServiceURLSet(t *testing.T) {
	serviceURL := urlx.MustParse("http://example.com")

	tests := []struct {
		name string
		env  Environment
		want bool
	}{
		{
			name: "URL is set",
			env: Environment{
				ServiceURL: *serviceURL,
			},
			want: true,
		},
		{
			name: "URL is not set",
			env:  Environment{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.env.IsServiceURLSet())
		})
	}
}

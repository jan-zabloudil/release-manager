package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_NewEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		creation EnvironmentCreation
		wantErr  bool
	}{
		{
			name: "Valid Environment",
			creation: EnvironmentCreation{
				ProjectID:     uuid.New(),
				Name:          "dev",
				ServiceRawURL: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Invalid Environment - not absolute service url",
			creation: EnvironmentCreation{
				ProjectID:     uuid.New(),
				Name:          "dev",
				ServiceRawURL: "example.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Environment - empty name",
			creation: EnvironmentCreation{
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

func TestEnvironment_toServiceURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{
			name:    "Valid Absolute URL",
			rawURL:  "http://example.com",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			rawURL:  "invalid",
			wantErr: true,
		},
		{
			name:    "Relative URL",
			rawURL:  "/relative/path",
			wantErr: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			_, err := toServiceURL(tt.rawURL)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnvironment_Update(t *testing.T) {
	validName := "New Name"
	validURL := "http://new.example.com"
	invalidName := ""
	invalidURL := "example"

	tests := []struct {
		name    string
		env     Environment
		update  EnvironmentUpdate
		wantErr bool
	}{
		{
			name: "Valid Update",
			env: Environment{
				Name: "Old Name",
			},
			update: EnvironmentUpdate{
				Name:          &validName,
				ServiceRawURL: &validURL,
			},
			wantErr: false,
		},
		{
			name: "Invalid Update - not absolute service url",
			env: Environment{
				Name: "Old Name",
			},
			update: EnvironmentUpdate{
				Name:          &validName,
				ServiceRawURL: &invalidURL,
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - empty name",
			env: Environment{
				Name: "Old Name",
			},
			update: EnvironmentUpdate{
				Name:          &invalidName,
				ServiceRawURL: &validURL,
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

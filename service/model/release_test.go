package model

import (
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRelease_NewRelease(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateReleaseInput
		wantErr bool
	}{
		{
			name: "Valid Release",
			input: CreateReleaseInput{
				ReleaseTitle: "Release 1.0",
				ReleaseNotes: "Initial release",
				GitTagName:   "v1.0.0",
				GitTagURL: url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/owner/repo/releases/tag/v1.0.0",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Release - Empty ReleaseTitle",
			input: CreateReleaseInput{
				ReleaseTitle: "",
				ReleaseNotes: "Initial release",
				GitTagName:   "v1.0.0",
				GitTagURL: url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/owner/repo/releases/tag/v1.0.0",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Release - missing git tag",
			input: CreateReleaseInput{
				ReleaseTitle: "Release 1.0",
				ReleaseNotes: "Initial release",
				GitTagName:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRelease(tt.input, uuid.New(), uuid.New())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRelease_Update(t *testing.T) {
	validTitle := "Updated Title"
	validNotes := "Updated Notes"
	InvalidTitle := ""

	tests := []struct {
		name    string
		input   UpdateReleaseInput
		want    Release
		wantErr bool
	}{
		{
			name: "Valid update",
			input: UpdateReleaseInput{
				ReleaseTitle: &validTitle,
				ReleaseNotes: &validNotes,
			},
			want: Release{
				ReleaseTitle: validTitle,
				ReleaseNotes: validNotes,
				GitTagName:   "v1.0.0",
				GitTagURL: url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/owner/repo/releases/tag/v1.0.0",
				},
			},
			wantErr: false,
		},
		{
			name: "No ReleaseTitle provided",
			input: UpdateReleaseInput{
				ReleaseTitle: &InvalidTitle,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Release{
				ID:           uuid.New(),
				ProjectID:    uuid.New(),
				ReleaseTitle: "Initial Title",
				ReleaseNotes: "Initial Notes",
				GitTagName:   "v1.0.0",
				GitTagURL: url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/owner/repo/releases/tag/v1.0.0",
				},
				AuthorUserID: uuid.New(),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			err := r.Update(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ReleaseTitle, r.ReleaseTitle)
				assert.Equal(t, tt.want.ReleaseNotes, r.ReleaseNotes)
				assert.True(t, r.UpdatedAt.After(tt.want.UpdatedAt))
			}
		})
	}
}

func TestCreateReleaseInput_ValidateGitTagName(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateReleaseInput
		wantErr bool
	}{
		{
			name: "Valid Git Tag Name",
			input: CreateReleaseInput{
				GitTagName: "v1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Empty Git Tag Name",
			input: CreateReleaseInput{
				GitTagName: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.ValidateGitTagName()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

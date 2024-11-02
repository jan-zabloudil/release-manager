package model

import (
	"net/url"
	"testing"
	"time"

	"release-manager/pkg/id"
	"release-manager/pkg/pointer"
	"release-manager/pkg/urlx"

	"github.com/stretchr/testify/assert"
)

func TestRelease_NewRelease(t *testing.T) {
	validURL := urlx.MustParse("https://github.com/owner/repo/releases/tag/v1.0.0")

	tests := []struct {
		name      string
		input     CreateReleaseInput
		gitTagURL url.URL
		wantErr   bool
	}{
		{
			name: "Valid Release",
			input: CreateReleaseInput{
				ReleaseTitle: "Release 1.0",
				ReleaseNotes: "Initial release",
				GitTagName:   "v1.0.0",
			},
			gitTagURL: *validURL,
			wantErr:   false,
		},
		{
			name: "Invalid Release - Empty ReleaseTitle",
			input: CreateReleaseInput{
				ReleaseTitle: "",
				ReleaseNotes: "Initial release",
				GitTagName:   "v1.0.0",
			},
			gitTagURL: *validURL,
			wantErr:   true,
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
		{
			name: "Invalid Release - empty URL",
			input: CreateReleaseInput{
				ReleaseTitle: "Release 1.0",
				ReleaseNotes: "Initial release",
				GitTagName:   "v1.0.0",
			},
			gitTagURL: url.URL{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRelease(tt.input, tt.gitTagURL, id.NewProject(), id.AuthUser{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRelease_Update(t *testing.T) {
	tests := []struct {
		name    string
		input   UpdateReleaseInput
		want    Release
		wantErr bool
	}{
		{
			name: "Valid update",
			input: UpdateReleaseInput{
				ReleaseTitle: pointer.StringPtr("Valid title"),
				ReleaseNotes: pointer.StringPtr("Valid notes"),
			},
			want: Release{
				ReleaseTitle: "Valid title",
				ReleaseNotes: "Valid notes",
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
				ReleaseTitle: pointer.StringPtr(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Release{
				ID:           id.NewRelease(),
				ProjectID:    id.NewProject(),
				ReleaseTitle: "Initial Title",
				ReleaseNotes: "Initial Notes",
				GitTagName:   "v1.0.0",
				GitTagURL: url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/owner/repo/releases/tag/v1.0.0",
				},
				AuthorUserID: id.AuthUser{},
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

func TestCreateReleaseInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateReleaseInput
		wantErr bool
	}{
		{
			name: "Valid input",
			input: CreateReleaseInput{
				ReleaseTitle: "Release 1.0",
				GitTagName:   "v1.0.0",
			},
			wantErr: false,
		},
		{
			name: "Invalid input - empty ReleaseTitle",
			input: CreateReleaseInput{
				ReleaseTitle: "",
				GitTagName:   "v1.0.0",
			},
			wantErr: true,
		},
		{
			name: "Invalid input - empty git tag name",
			input: CreateReleaseInput{
				ReleaseTitle: "New release",
				GitTagName:   "",
			},
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

func TestGithubGeneratedReleaseNotesInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   GithubReleaseNotesInput
		wantErr bool
	}{
		{
			name: "Valid GitTagName",
			input: GithubReleaseNotesInput{
				GitTagName: pointer.StringPtr("v1.0.0"),
			},
			wantErr: false,
		},
		{
			name: "Nil GitTagName",
			input: GithubReleaseNotesInput{
				GitTagName: nil,
			},
			wantErr: true,
		},
		{
			name: "Empty GitTagName",
			input: GithubReleaseNotesInput{
				GitTagName: pointer.StringPtr(""),
			},
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

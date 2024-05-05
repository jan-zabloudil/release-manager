package model

import (
	"testing"

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
			},
			wantErr: false,
		},
		{
			name: "Invalid Release - Empty ReleaseTitle",
			input: CreateReleaseInput{
				ReleaseTitle: "",
				ReleaseNotes: "Initial release",
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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_NewEnvURL(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid absolute URL",
			input:   "https://example.com",
			wantErr: false,
		},
		{
			name:    "Valid absolute URL with path",
			input:   "https://example.com/path",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			input:   "://example.com",
			wantErr: true,
		},
		{
			name:    "Invalid URL without scheme",
			input:   "example.com",
			wantErr: true,
		},
		{
			name:    "Relative URL",
			input:   "/path",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewEnvURL(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

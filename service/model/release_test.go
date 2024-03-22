package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSourceCode(t *testing.T) {
	testCases := []struct {
		name            string
		tag             string
		targetCommitIsh string
		wantErr         bool
	}{
		{
			name:            "Empty tag",
			tag:             "",
			targetCommitIsh: "targetCommitIsh",
			wantErr:         true,
		},
		{
			name:            "Valid parameters",
			tag:             "tag",
			targetCommitIsh: "targetCommitIsh",
			wantErr:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSourceCode(tc.tag, tc.targetCommitIsh)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

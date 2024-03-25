package model

import (
	"testing"

	svcerr "release-manager/service/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplateType(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		expectedErr error
	}{
		{
			name:        "valid template type",
			key:         releaseMsgTmplType,
			expectedErr: nil,
		},
		{
			name:        "invalid template type",
			key:         "invalidType",
			expectedErr: svcerr.ErrInvalidTemplateType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTemplateType(tt.key)

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

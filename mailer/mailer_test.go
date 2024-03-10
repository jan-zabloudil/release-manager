package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEmailRequest(t *testing.T) {
	tests := []struct {
		name              string
		testingMode       bool
		testRecipient     string
		realRecipient     []string
		expectedRecipient []string
	}{
		{
			name:              "email request with real recipients",
			testingMode:       false,
			testRecipient:     "test@test.tt",
			realRecipient:     []string{"prod@test.tt"},
			expectedRecipient: []string{"prod@test.tt"},
		},
		{
			name:              "email request with test recipient",
			testingMode:       true,
			testRecipient:     "test@test.tt",
			realRecipient:     []string{"prod@test.tt"},
			expectedRecipient: []string{"test@test.tt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				TestingMode:   tt.testingMode,
				TestRecipient: tt.testRecipient,
			}
			mailer := New(cfg)

			result, _ := mailer.buildEmailRequest(tt.realRecipient, "test.html")

			assert.Equal(t, tt.expectedRecipient, result.To, "wrong recipient")
		})
	}
}

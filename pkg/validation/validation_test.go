package validation_test

import (
	"testing"

	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateURL(t *testing.T) {
	t.Parallel()

	testValidator := validation.New()

	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "url is valid with github host",
			input:       "git@github.com:VladPetriv/nvim-config.git",
			expectedErr: nil,
		},
		{
			name:        "url is valid with gitlab host",
			input:       "git@gitlab.com:gitlab-org/gitaly.git",
			expectedErr: nil,
		},
		{
			name:        "url is not valid",
			input:       "test",
			expectedErr: validation.ErrURLContainsInvalidHost,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := testValidator.ValidateURL(tt.input)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

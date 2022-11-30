package validation_test

import (
	"fmt"
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

func Test_ValidateRepoFiles(t *testing.T) {
	t.Parallel()

	testValidator := validation.New()

	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "repository contains valid files",
			input:       "../../../../../.config/nvim/",
			expectedErr: nil,
		},
		{
			name:        "nvim is not a main directory",
			input:       "../../../../../.config/",
			expectedErr: fmt.Errorf("failed to get list of repository files: %w", validation.ErrNvimConfigDirIsNotMainDir),
		},
		{
			name:        "repository didn't contains valid files",
			input:       "../",
			expectedErr: validation.ErrNoBaseFilesInRepository,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := testValidator.ValidateRepoFiles(tt.input)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

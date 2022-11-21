package validation_test

import (
	"fmt"
	"testing"

	"github.com/VladPetriv/setup-neovim/pkg/logger"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
)

var testValidator = validation.New(logger.Get())

func Test_ValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "url is valid with github",
			input:       "git@github.com:VladPetriv/nvim-config.git",
			expectedErr: nil,
		},
		{
			name:        "url is valid with gitlab",
			input:       "git@gitlab.com:gitlab-org/gitaly.git",
			expectedErr: nil,
		},
		{
			name:        "url is not valid",
			input:       "test",
			expectedErr: fmt.Errorf("url contains invalid host: test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testValidator.ValidateURL(tt.input)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Test_ValidateRepoFiles(t *testing.T) {
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
			name:        "repository didn't contains valid files",
			input:       "../",
			expectedErr: fmt.Errorf("repository didn't contains base files for neovim configuration"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testValidator.ValidateRepoFiles(tt.input)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

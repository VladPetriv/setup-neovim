package validation_test

import (
	"os"
	"testing"

	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateURL(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := testValidator.ValidateURL(tt.input)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidator_ValidateRepoFiles(t *testing.T) { //nolint:tparallel,lll // t.Parallel() causes conflicts with directories
	t.Parallel()

	testValidator := validation.New()

	tests := []struct {
		name         string
		input        string
		withBaseFile bool
		expectedErr  error
	}{
		{
			name:         "validation for repository files successfully completed",
			input:        "./nvim",
			withBaseFile: true,
			expectedErr:  nil,
		},
		{
			name:        "validation for repository files failed with no required files",
			input:       "./nvim",
			expectedErr: validation.ErrNoBaseFilesInRepository,
		},
		{
			name:        "validation for repository files failed with not found path",
			input:       "./nvim_test",
			expectedErr: validation.ErrDirectoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Mkdir("nvim", 0o755)
			if err != nil {
				t.Fatal(err)
			}

			if tt.withBaseFile {
				_, err = os.Create("./nvim/init.lua")
				if err != nil {
					t.Fatal(err)
				}
			}

			err = testValidator.ValidateRepoFiles(tt.input)
			assert.Equal(t, tt.expectedErr, err)

			err = os.RemoveAll("./nvim")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

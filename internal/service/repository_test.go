package service_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CloneAndValidateRepository(t *testing.T) { //nolint:tparallel,lll // t.Parallel() causes conflicts with go-git
	t.Parallel()

	testService := service.NewRepository(input.New(), validation.New())

	tests := []struct {
		name          string
		inputURL      string
		inputFilePath string
		wantErr       bool
	}{
		{
			name:          "CloneAndValidateRepository success with https URL [github]",
			inputURL:      "https://github.com/jdhao/nvim-config.git",
			inputFilePath: "",
		},
		{
			name:          "CloneAndValidateRepository success with https URL [gitlab]",
			inputURL:      "https://gitlab.com/hantamkoding-dotfiles/neovim.git",
			inputFilePath: "",
		},
		{
			name:          "CloneAndValidateRepository success with SSH URL",
			inputURL:      "git@github.com:VladPetriv/nvim-config.git",
			inputFilePath: ".ssh/id_ed25519",
		},
		{
			name:          "CloneAndValidateRepository fail with file validation",
			inputURL:      "https://gitlab.com/VladPetriv/tg_scanner.git",
			inputFilePath: "",
			wantErr:       true,
		},
		{
			name:          "CloneAndValidateRepository fail with error when create public key from file",
			inputURL:      "git@github.com:VladPetriv/nvim-config.git",
			inputFilePath: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdin bytes.Buffer
			stdin.Write([]byte(fmt.Sprintf("%s\n", tt.inputFilePath)))

			err := testService.CloneAndValidateRepository(tt.inputURL, &stdin)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			err = os.RemoveAll("./nvim")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRepository_ProcessUserURL(t *testing.T) {
	t.Parallel()

	testService := service.NewRepository(input.New(), validation.New())

	tests := []struct {
		name        string
		input       string
		expectedURL string
		expectedErr error
	}{
		{
			name:        "ProcessUserURL success with valid host[github]",
			input:       "git@github.com:VladPetriv/setup-neovim.git",
			expectedURL: "git@github.com:VladPetriv/setup-neovim.git",
			expectedErr: nil,
		},
		{
			name:        "ProcessUserURL success with valid host[gitlab]",
			input:       "git@gitlab.com:gitlab-org/gitaly.git",
			expectedURL: "git@gitlab.com:gitlab-org/gitaly.git",
			expectedErr: nil,
		},
		{
			name:        "ProcessUserURL fail with invalid host",
			input:       "test",
			expectedURL: "",
			expectedErr: validation.ErrURLContainsInvalidHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdin bytes.Buffer
			stdin.Write([]byte(fmt.Sprintf("%s\n", tt.input)))

			url, err := testService.ProcessUserURL(&stdin)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedURL, url)
		})
	}
}

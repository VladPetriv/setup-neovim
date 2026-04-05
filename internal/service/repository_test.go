package service_test

import (
	"os"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CloneRepository(t *testing.T) { //nolint:tparallel // t.Parallel() causes conflicts with go-git
	t.Parallel()

	testService := service.NewRepository(validation.New())

	tests := []struct {
		name       string
		url        string
		sshKeyPath string
		wantErr    bool
	}{
		{
			name: "success with HTTPS URL [github]",
			url:  "https://github.com/jdhao/nvim-config.git",
		},
		{
			name: "success with HTTPS URL [gitlab]",
			url:  "https://gitlab.com/hantamkoding-dotfiles/neovim.git",
		},
		{
			name:       "success with SSH URL",
			url:        "git@github.com:VladPetriv/nvim-config.git",
			sshKeyPath: ".ssh/id_ed25519",
		},
		{
			name:    "fails when SSH key file is missing",
			url:     "git@github.com:VladPetriv/nvim-config.git",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testService.CloneRepository(tt.url, tt.sshKeyPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			os.RemoveAll("./nvim")
		})
	}
}

func TestRepository_ValidateRepository(t *testing.T) {
	t.Parallel()

	testService := service.NewRepository(validation.New())

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "fails when directory does not exist",
			path:    "./nonexistent_repo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := testService.ValidateRepository(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

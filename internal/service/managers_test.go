package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_DetectInstalledPackageManagers(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	testService := service.NewManager()

	type precondition struct {
		createPackerDir  bool
		createVimPlugDir bool
	}

	tests := []struct {
		name         string
		precondition precondition
		wantCount    int
	}{
		{
			name: "detects both package managers",
			precondition: precondition{
				createPackerDir:  true,
				createVimPlugDir: true,
			},
			wantCount: 2,
		},
		{
			name: "detects one package manager",
			precondition: precondition{
				createPackerDir: true,
			},
			wantCount: 1,
		},
		{
			name:      "detects no package managers",
			wantCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir))
			os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir))

			if tt.precondition.createPackerDir {
				err = os.MkdirAll(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir), 0o755)
				require.NoError(t, err)
			}
			if tt.precondition.createVimPlugDir {
				err = os.MkdirAll(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir), 0o755)
				require.NoError(t, err)
			}

			t.Cleanup(func() {
				os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir))
				os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir))
			})

			var detected []string
			detected, err = testService.DetectInstalledPackageManagers()
			assert.NoError(t, err)
			assert.Len(t, detected, tt.wantCount)
		})
	}
}

func TestManager_DeletePackageManagers(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	testService := service.NewManager()

	packerDir := fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir)
	vimPlugDir := fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir)

	err = os.MkdirAll(packerDir, 0o755)
	require.NoError(t, err)
	err = os.MkdirAll(vimPlugDir, 0o755)
	require.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(packerDir)
		os.RemoveAll(vimPlugDir)
	})

	err = testService.DeletePackageManagers()
	assert.NoError(t, err)

	_, err = os.Lstat(packerDir)
	assert.Error(t, err, "packer directory should have been deleted")

	_, err = os.Lstat(vimPlugDir)
	assert.Error(t, err, "vim-plug directory should have been deleted")
}

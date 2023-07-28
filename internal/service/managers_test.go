package service_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_DetectInstalledPackageManagers(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	testService := service.New(&service.Options{
		Inputter:  input.New(),
		Validator: validation.New(),
	})

	type precondition struct {
		createPackerDir  bool
		createVimPlugDir bool
	}

	type expected struct {
		count int
	}

	tests := []struct {
		name         string
		precondition precondition
		expected     expected
	}{
		{
			name: "DetectInstalledPackageManagers successfully with 2 detected package managers",
			precondition: precondition{
				createPackerDir:  true,
				createVimPlugDir: true,
			},
			expected: expected{
				count: 2,
			},
		},
		{
			name: "DetectInstalledPackageManagers successfully with 1 detected package managers",
			precondition: precondition{
				createPackerDir: true,
			},
			expected: expected{
				count: 1,
			},
		},
		{
			name:     "DetectInstalledPackageManagers failed with no detected package managers",
			expected: expected{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err = os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir))
			assert.NoError(t, err)
			err = os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir))
			assert.NoError(t, err)

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
				assert.NoError(t, err)

				err = os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir))
				assert.NoError(t, err)
			})

			_, actualCount, actualErr := testService.DetectInstalledPackageManagers()
			assert.NoError(t, actualErr)
			assert.Equal(t, tt.expected.count, actualCount)
		})
	}
}

func TestManager_ProcessInputForPackageManagers(t *testing.T) {
	t.Parallel()

	testService := service.New(&service.Options{
		Inputter:  input.New(),
		Validator: validation.New(),
	})

	tests := []struct {
		name                   string
		inputHaveInstalled     string
		inputPackageManager    string
		expectedPackageManager string
	}{
		{
			name:                   "ProcessInputForPackageManagers with no needs to package manager",
			inputHaveInstalled:     "y",
			expectedPackageManager: "skip",
		},
		{
			name:                   "ProcessInputForPackageManagers with packer",
			inputHaveInstalled:     "n",
			inputPackageManager:    "packer",
			expectedPackageManager: "packer",
		},
		{
			name:                   "ProcessInputForPackageManagers with vim-plug",
			inputHaveInstalled:     "n",
			inputPackageManager:    "vim-plug",
			expectedPackageManager: "vim-plug",
		},
		{
			name:                   "ProcessInputForPackageManagers with empty result",
			inputHaveInstalled:     "n",
			inputPackageManager:    "",
			expectedPackageManager: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := strings.NewReader(
				fmt.Sprintf("%s\n%s\n", tt.inputHaveInstalled, tt.inputPackageManager),
			)

			got, err := testService.GetPackageMangerIfNotInstalled(input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPackageManager, got)
		})
	}
}

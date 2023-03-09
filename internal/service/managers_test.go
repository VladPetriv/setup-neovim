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
)

func Test_ProcessPackageManagers(t *testing.T) { //nolint:tparallel // t.Parallel() causes conflicts with files
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
		wantErr                bool
		expectedError          error
	}{
		{
			name:                   "ProcessPackageManagers with no needs to install package manager",
			inputHaveInstalled:     "y",
			expectedPackageManager: "",
		},
		{
			name:                   "ProcessPackageManagers with installing packer",
			inputHaveInstalled:     "n",
			inputPackageManager:    "packer",
			expectedPackageManager: "packer",
		},
		{
			name:                   "ProcessPackageManagers with installing vim plug",
			inputHaveInstalled:     "n",
			inputPackageManager:    "vim-plug",
			expectedPackageManager: "vim-plug",
		},
		{
			name:               "ProcessPackageManagers with invalid input",
			inputHaveInstalled: "",
			wantErr:            true,
			expectedError:      service.ErrEnterValidAnswer,
		},

		{
			name:                "ProcessPackageManagers with invalid input",
			inputHaveInstalled:  "n",
			inputPackageManager: "hello",
			wantErr:             true,
			expectedError:       service.ErrEnterValidAnswer,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(
				fmt.Sprintf("%s\n%s\n", tt.inputHaveInstalled, tt.inputPackageManager),
			)

			got, err := testService.ProcessPackageManagers(input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPackageManager, got)
			}

			if tt.expectedPackageManager == "packer" {
				homeDir, _ := os.UserHomeDir()

				err = os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/", homeDir))
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func Test_ProcessInputForPackageManagers(t *testing.T) {
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

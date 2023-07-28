package service_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
)

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

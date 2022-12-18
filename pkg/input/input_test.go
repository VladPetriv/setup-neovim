package input_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/stretchr/testify/assert"
)

func Test_GetInput(t *testing.T) {
	t.Parallel()

	testInput := input.New()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "GetInput success",
			input: "test",
			want:  "test",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdin bytes.Buffer
			stdin.Write([]byte(fmt.Sprintf("%s\n", tt.input)))

			got, err := testInput.GetInput(&stdin)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, "", got)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func Test_ProcessInputForPackageManagers(t *testing.T) {
	t.Parallel()

	testInput := input.New()

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

			got, err := testInput.ProcessInputForPackageManagers(input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPackageManager, got)
		})
	}
}

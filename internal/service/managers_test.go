package service_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_DetectInstalledPackageManagers(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	testService := service.NewManager(input.New())

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

func TestManager_ProcessAlreadyInstalledPackageManagers(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	testService := service.NewManager(input.New())

	type precondition struct {
		createPackerDir  bool
		createVimPlugDir bool
	}

	type args struct {
		count           int
		wantRemoveInput string
	}

	type expected struct {
		shouldInstall bool
		err           error
	}

	tests := []struct {
		name         string
		precondition precondition
		args         args
		expected     expected
	}{
		{
			name: "ProcessAlreadyInstalledPackageManagers successful with 2 deleted managers",
			precondition: precondition{
				createPackerDir:  true,
				createVimPlugDir: true,
			},
			args: args{
				count:           2,
				wantRemoveInput: "y",
			},
			expected: expected{
				shouldInstall: true,
			},
		},
		{
			name: "ProcessAlreadyInstalledPackageManagers successful with 1 deleted managers",
			precondition: precondition{
				createPackerDir: true,
			},
			args: args{
				count:           1,
				wantRemoveInput: "y",
			},
			expected: expected{
				shouldInstall: true,
			},
		},
		{
			name: "ProcessAlreadyInstalledPackageManagers successful with no detected managers",
			precondition: precondition{
				createPackerDir: true,
			},
			args: args{
				count: 0,
			},
			expected: expected{
				shouldInstall: true,
			},
		},
		{
			name: "ProcessAlreadyInstalledPackageManagers failed with no need to delete",
			precondition: precondition{
				createPackerDir: true,
			},
			args: args{
				count:           1,
				wantRemoveInput: "n",
			},
			expected: expected{
				shouldInstall: false,
				err:           nil,
			},
		},
		{
			name: "ProcessAlreadyInstalledPackageManagers failed with invalid answer",
			precondition: precondition{
				createPackerDir: true,
			},
			args: args{
				count:           1,
				wantRemoveInput: "invalid",
			},
			expected: expected{
				err: service.ErrEnterValidAnswer,
			},
		},
	}
	for _, tt := range tests {
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

			input := strings.NewReader(
				fmt.Sprintf("%s\n", tt.args.wantRemoveInput),
			)

			actual, processErr := testService.ProcessAlreadyInstalledPackageManagers(tt.args.count, input)
			assert.Equal(t, tt.expected.err, processErr)
			assert.Equal(t, tt.expected.shouldInstall, actual)
		})
	}
}

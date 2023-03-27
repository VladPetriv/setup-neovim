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

func TestFile_ExtractAndMoveConfigDirectory(t *testing.T) { //nolint:tparallel,lll // t.Parallel() causes conflicts with dirs
	t.Parallel()

	testService := service.New(&service.Options{
		Inputter:  input.New(),
		Validator: validation.New(),
	})

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		input         string
		wantErr       bool
		directoryType string
	}{
		{
			name:    "failed by directory not found",
			input:   "./test_not_found",
			wantErr: true,
		},
		{
			name:          "successful by moving config as main directory",
			input:         "./nvim",
			directoryType: "main",
		},
		{
			name:          "successful by moving config as extracted directory",
			input:         "./test_not_main",
			directoryType: "not_main",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err = createDirectoryByType(tt.directoryType)
			require.NoError(t, err)

			got := testService.ExtractAndMoveConfigDirectory(tt.input)
			if tt.wantErr {
				assert.Error(t, got)
			} else {
				assert.NoError(t, got)
			}

			err = os.RemoveAll(fmt.Sprintf("%s/.config/nvim", homeDir))
			if err != nil {
				require.NoError(t, err)
			}
		})
	}
}

func createDirectoryByType(directoryType string) error {
	directories := map[string]string{
		"main":     "./nvim",
		"not_main": "./test_not_main/nvim",
	}

	var directory string

	switch directoryType {
	case "main":
		directory = directories["main"]
	case "not_main":
		directory = directories["not_main"]
	default:
		return nil
	}

	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	return nil
}

func TestFile_DeleteConfigOrStopInstallation(t *testing.T) { //nolint:tparallel,lll // t.Parallel() causes conflicts with dirs
	t.Parallel()

	testService := service.New(&service.Options{
		Inputter:  input.New(),
		Validator: validation.New(),
	})

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	configPath := fmt.Sprintf("%s/.config/nvim", homeDir)

	tests := []struct {
		name                        string
		input                       string
		shouldCreateConfigDirectory bool
		want                        error
	}{
		{
			name:                        "DeleteConfigOrStopInstallation successful",
			input:                       "y",
			want:                        nil,
			shouldCreateConfigDirectory: true,
		},
		{
			name:  "DeleteConfigOrStopInstallation successful with no needs to deleting config",
			input: "",
			want:  nil,
		},
		{
			name:                        "DeleteConfigOrStopInstallation failed with stop installation",
			input:                       "n",
			want:                        service.ErrStopInstallation,
			shouldCreateConfigDirectory: true,
		},
		{
			name:                        "DeleteConfigOrStopInstallation failed with invalid answer",
			input:                       "hello",
			want:                        service.ErrEnterValidAnswer,
			shouldCreateConfigDirectory: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(
				fmt.Sprintf("%s\n", tt.input),
			)

			if tt.shouldCreateConfigDirectory {
				err = os.MkdirAll(configPath, 0755)
				require.NoError(t, err)
			}

			got := testService.DeleteConfigOrStopInstallation(input)
			assert.Equal(t, tt.want, got)

			if tt.shouldCreateConfigDirectory {
				err = os.RemoveAll(configPath)
				require.NoError(t, err)
			}
		})
	}
}

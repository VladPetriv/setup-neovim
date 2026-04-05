package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFile_ExtractAndMoveConfigDirectory(t *testing.T) {
	testService := service.NewFile(validation.New())

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
		t.Run(tt.name, func(t *testing.T) {
			err = createDirectoryByType(tt.directoryType)
			assert.NoError(t, err)

			t.Cleanup(func() {
				err = os.RemoveAll(fmt.Sprintf("%s/.config/nvim", homeDir))
				if err != nil {
					require.NoError(t, err)
				}
			})

			got := testService.ExtractAndMoveConfigDirectory(tt.input)
			if tt.wantErr {
				assert.Error(t, got)
			} else {
				assert.NoError(t, got)
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

func TestFile_CheckConfigExists(t *testing.T) { //nolint:tparallel // t.Parallel() causes conflicts with dirs
	t.Parallel()

	testService := service.NewFile(validation.New())

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	configPath := fmt.Sprintf("%s/.config/nvim", homeDir)

	err = os.RemoveAll(configPath)
	require.NoError(t, err)

	tests := []struct {
		name         string
		createConfig bool
		wantExists   bool
	}{
		{
			name:         "returns true when config directory exists",
			createConfig: true,
			wantExists:   true,
		},
		{
			name:       "returns false when config directory is absent",
			wantExists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createConfig {
				err = os.MkdirAll(configPath, 0755)
				require.NoError(t, err)
			}

			t.Cleanup(func() {
				os.RemoveAll(configPath)
			})

			var got bool
			got, err = testService.CheckConfigExists()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantExists, got)
		})
	}
}

func TestFile_DeleteConfig(t *testing.T) {
	t.Parallel()

	testService := service.NewFile(validation.New())

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	configPath := fmt.Sprintf("%s/.config/nvim", homeDir)

	err = os.MkdirAll(configPath, 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(configPath)
	})

	err = testService.DeleteConfig()
	assert.NoError(t, err)

	_, err = os.Lstat(configPath)
	assert.Error(t, err, "config directory should have been deleted")
}

func TestFile_DeleteRepositoryDirectory(t *testing.T) {
	testService := service.NewFile(validation.New())

	tests := []struct {
		name            string
		input           string
		createDirectory bool
	}{
		{
			name:  "does nothing when path is empty",
			input: "",
		},
		{
			name:            "successfully deletes existing directory",
			input:           "./test_delete_found",
			createDirectory: true,
		},
		{
			name:  "does nothing when directory does not exist",
			input: "./test_delete_not_found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createDirectory {
				err := os.MkdirAll(tt.input, 0755)
				require.NoError(t, err)
			}

			got := testService.DeleteRepositoryDirectory(tt.input)
			assert.NoError(t, got)

			if tt.createDirectory {
				_, err := os.Lstat(tt.input)
				assert.Errorf(t, err, "input directory is not deleted. Directory path: %s", tt.input)
			}
		})
	}
}

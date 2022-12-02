package service

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var ErrExtractDir = errors.New("failed to")

func (s service) MoveConfigDirectory() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get path to the home directory: %w", err)
	}

	err = os.Rename("./nvim", fmt.Sprintf("%s/.config/nvim", homeDir))
	if err != nil {
		removeErr := os.RemoveAll("./nvim")
		if removeErr != nil {
			return fmt.Errorf("moving config failed, failed to remove repository: %w", err)
		}

		return fmt.Errorf("failed to move repository into .config directory: %w", err)
	}

	return nil
}

func (s service) ExtractConfigDirectory(repositoryPath string) (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("failed to get path to the home directory: %w", err)
	}

	configPath, err := getConfigPath(repositoryPath)
	if err != nil {
		return false, fmt.Errorf("failed to extract config directory: %w", err)
	}

	if configPath != "" {
		err = os.Rename(configPath, fmt.Sprintf("%s/.config/nvim", homeDir))
		if err != nil {
			return false, fmt.Errorf("failed to move config directory: %w", err)
		}

		return true, nil
	}

	return false, nil
}

func getConfigPath(repositoryPath string) (string, error) {
	var dirPath string
	var filesCount int //nolint

	err := filepath.Walk(repositoryPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk in repository error: %w", err)
		}

		if filesCount > 0 && info.Name() == "nvim" && info.IsDir() {
			dirPath = path

			return nil
		}

		filesCount++

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to get path to nvim config in repository: %w", err)
	}

	return dirPath, nil
}

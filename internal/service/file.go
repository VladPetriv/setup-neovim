package service

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func (s service) ExtractAndMoveConfigDirectory(repositoryPath string) error {
	configPath, err := getConfigPath(repositoryPath)
	if err != nil {
		return fmt.Errorf("failed to config path: %w", err)
	}

	err = moveConfigDirectory(configPath)
	if err != nil {
		if errors.Is(err, ErrDirectoryNotFound) {
			return err
		}

		return fmt.Errorf("failed to move config directory: %w", err)
	}

	return nil
}

func moveConfigDirectory(configPath string) error {
	if configPath == "" {
		if _, err := os.Lstat("./nvim"); err != nil {
			return ErrDirectoryNotFound
		}

		configPath = "./nvim"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get path to the home directory: %w", err)
	}

	err = os.Rename(configPath, fmt.Sprintf("%s/.config/nvim", homeDir))
	if err != nil {
		removeErr := os.RemoveAll(configPath)
		if removeErr != nil {
			return fmt.Errorf("moving config failed, failed to remove repository: %w", err)
		}

		return fmt.Errorf("failed to move repository into .config directory: %w", err)
	}

	return nil
}

// getConfigPath is used for get path to nvim config if it's not main directory.
func getConfigPath(repositoryPath string) (string, error) {
	var dirPath string
	var filesCount int

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

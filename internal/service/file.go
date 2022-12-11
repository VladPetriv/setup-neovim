package service

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func (s service) ExtractAndMoveConfigDirectory(repositoryPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get path to the home directory: %w", err)
	}

	configPath, err := getConfigPath(repositoryPath)
	if err != nil {
		return fmt.Errorf("failed to config path: %w", err)
	}

	// if config path is not empty, it means that nvim config was not main directory in the repository
	if configPath != "" {
		err = os.Rename(configPath, fmt.Sprintf("%s/.config/nvim", homeDir))
		if err != nil {
			removeErr := os.RemoveAll("./nvim")
			if removeErr != nil {
				return fmt.Errorf("moving config failed, failed to remove repository: %w", removeErr)
			}

			return fmt.Errorf("failed to move repository into .config directory: %w", err)
		}

		return nil
	}

	// if config path is empty, it means that nvim is main directory in the repository
	err = moveConfigDirectory()
	if err != nil {
		return fmt.Errorf("failed to move config directory: %w", err)
	}

	return nil
}

func moveConfigDirectory() error {
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

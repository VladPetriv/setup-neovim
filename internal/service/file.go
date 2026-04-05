package service

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

type fileService struct {
	validator validation.Validator
}

func NewFile(validator validation.Validator) FileService {
	return &fileService{validator: validator}
}

func (f fileService) CheckUtilStatus() map[string]string {
	return f.validator.ValidateConsoleUtilities()
}

func (f fileService) CheckConfigExists() (bool, error) {
	return checkIfConfigDirectoryIsExist()
}

func (f fileService) DeleteConfig() error {
	return deleteConfig()
}

func (f fileService) ExtractAndMoveConfigDirectory(repositoryPath string) error {
	configPath, err := getConfigPath(repositoryPath)
	if err != nil {
		if errors.Is(err, ErrDirectoryNotFound) {
			return err
		}

		return fmt.Errorf("get config path: %w", err)
	}

	err = moveConfigDirectory(configPath)
	if err != nil {
		return fmt.Errorf("move config directory: %w", err)
	}

	return nil
}

func (f fileService) DeleteRepositoryDirectory(path string) error {
	if path == "" {
		return nil
	}

	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("remove directory by path: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the nvim config within the cloned repository.
func getConfigPath(repositoryPath string) (string, error) {
	var dirPath string
	var filesCount int

	err := filepath.Walk(repositoryPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk through files in directory: %w", err)
		}

		if filesCount > 0 && info.Name() == "nvim" && info.IsDir() {
			dirPath = path

			return nil
		}

		filesCount++

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("get path to config inside repository: %w", err)
	}

	if dirPath == "" {
		defaultConfigPath := fmt.Sprintf("./%s", DirectoryNameForClonnedRepository)
		_, statErr := os.Lstat(defaultConfigPath)
		if statErr != nil {
			return "", ErrDirectoryNotFound
		}

		return defaultConfigPath, nil
	}

	return dirPath, nil
}

func moveConfigDirectory(configPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	err = os.Rename(configPath, fmt.Sprintf("%s/%s", homeDir, systemNeovimConfigPath))
	if err != nil {
		return fmt.Errorf("can't move repository into .config directory: %w", err)
	}

	return nil
}

func checkIfConfigDirectoryIsExist() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("get home directory: %w", err)
	}

	_, err = os.Lstat(fmt.Sprintf("%s/%s", homeDir, systemNeovimConfigPath))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("check config existence: %w", err)
	}

	return true, nil
}

func deleteConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	err = os.RemoveAll(fmt.Sprintf("%s/%s", homeDir, systemNeovimConfigPath))
	if err != nil {
		return fmt.Errorf("remove existing config directory: %w", err)
	}

	return nil
}

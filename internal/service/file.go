package service

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

type fileService struct {
	inputter  input.Inputter
	validator validation.Validator
}

func NewFile(inputter input.Inputter, validator validation.Validator) *fileService {
	return &fileService{
		inputter:  inputter,
		validator: validator,
	}
}

func (f fileService) CheckUtilStatus() map[string]string {
	return f.validator.ValidateConsoleUtilities()
}

func (f fileService) ExtractAndMoveConfigDirectory(repositoryPath string) error {
	configPath, err := getConfigPath(repositoryPath)
	if err != nil {
		return fmt.Errorf("get config path: %w", err)
	}

	err = moveConfigDirectory(configPath)
	if err != nil {
		if errors.Is(err, ErrDirectoryNotFound) {
			return err
		}

		return fmt.Errorf("move config directory: %w", err)
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
		return fmt.Errorf("get home directory: %w", err)
	}

	err = os.Rename(configPath, fmt.Sprintf("%s/.config/nvim", homeDir))
	if err != nil {
		removeErr := os.RemoveAll(configPath)
		if removeErr != nil {
			return fmt.Errorf("moving config failed, can't remove old config path: %w", err)
		}

		return fmt.Errorf("can't move repository into .config directory: %w", err)
	}

	return nil
}

// getConfigPath is used for get path to nvim config if it's not main directory.
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

	return dirPath, nil
}

func (f fileService) DeleteConfigOrStopInstallation(stdin io.Reader) error {
	exist, err := checkIfConfigDirectoryIsExist()
	if err != nil {
		// NOTE: When directory with config not found we need to continue installation process
		if !errors.Is(err, ErrDirectoryNotFound) {
			return fmt.Errorf("check if config directory already exist: %w", err)
		}
	}

	// directory not found we should continue installation
	if !exist {
		return ErrConfigNotFound
	}

	fmt.Printf("Already installed neovim config detected...\nDo you want to remove it for continue installation?(y/n):")
	shouldStopOrContinueInstallation, err := f.inputter.GetInput(stdin)
	if err != nil {
		return fmt.Errorf("get user input: %w", err)
	}

	switch shouldStopOrContinueInstallation {
	case "y":
		return deleteConfig()
	case "n":
		return ErrStopInstallation
	default:
		return ErrEnterValidAnswer
	}
}

// checkIfConfigDirectoryIsExist checks if config directory existed by default path.
func checkIfConfigDirectoryIsExist() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("get home directory: %w", err)
	}

	if _, err = os.Lstat(fmt.Sprintf("%s/.config/nvim", homeDir)); err != nil {
		return false, ErrDirectoryNotFound
	}

	return true, nil
}

func deleteConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	if err = os.RemoveAll(fmt.Sprintf("%s/.config/nvim", homeDir)); err != nil {
		return fmt.Errorf("remove existed config directory: %w", err)
	}

	return nil
}

package validation

import (
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrURLContainsInvalidHost    = errors.New("url contains invalid host")
	ErrNoBaseFilesInRepository   = errors.New("repository didn't contains base files for nvim configuration")
	ErrNvimConfigDirIsNotMainDir = errors.New("directory with nvim config is not the main directory")
)

type validation struct{}

var _ Validator = (*validation)(nil)

func New() *validation { //nolint
	return &validation{}
}

func (v validation) ValidateURL(url string) error {
	var errCount int

	hosts := [2]string{"github.com", "gitlab.com"}
	for _, host := range hosts {
		if strings.Contains(url, host) {
			continue
		}

		errCount++
	}

	if errCount > 1 {
		return ErrURLContainsInvalidHost
	}

	return nil
}

func (v validation) ValidateRepoFiles(path string) error {
	var errCount int

	files, err := getRepositoryFiles(path)
	if err != nil {
		return fmt.Errorf("failed to get list of repository files: %w", err)
	}

	baseFiles := [2]string{"init.lua", "init.vim"}
	for _, file := range baseFiles {
		if strings.Contains(files, file) {
			continue
		}

		errCount++
	}

	if errCount > 1 {
		return ErrNoBaseFilesInRepository
	}

	return nil
}

func getRepositoryFiles(path string) (string, error) {
	var files string
	var filesCount int

	err := filepath.Walk(path, func(_ string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk in repository error: %w", err)
		}

		if filesCount > 0 && info.Name() == "nvim" && info.IsDir() {
			return ErrNvimConfigDirIsNotMainDir
		}

		files += fmt.Sprintf(" %s", info.Name())
		filesCount++

		return nil
	})
	if err != nil {
		if errors.Is(err, ErrNvimConfigDirIsNotMainDir) {
			return "", err
		}

		return "", fmt.Errorf("get list of files error: %w", err)
	}

	return files, nil
}

func (v validation) ValidateConsoleUtilities() map[string]string {
	errors := make(map[string]string)

	utils := [2]string{"nvim", "git"}
	for _, util := range utils {
		cmd := exec.Command(util, "--version")

		if err := cmd.Run(); err != nil {
			errors[util] = fmt.Sprintf("Please install %s before using setup nvim util", util)
		}
	}

	return errors
}

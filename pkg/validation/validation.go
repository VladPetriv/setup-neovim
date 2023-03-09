package validation

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrURLContainsInvalidHost  = errors.New("url contains invalid host")
	ErrNoBaseFilesInRepository = errors.New("repository didn't contains base files for nvim configuration")
	ErrDirectoryNotFound       = errors.New("directory not found")
)

type validation struct{}

var _ Validator = (*validation)(nil)

func New() Validator {
	return &validation{}
}

func (v validation) ValidateURL(url string) error {
	// allowableErrorCount represents the number of allowed errors.
	allowableErrorCount := 2

	var errCount int

	hosts := [3]string{"github.com", "gitlab.com", "bitbucket.org"}
	for _, host := range hosts {
		if strings.Contains(url, host) {
			continue
		}

		errCount++
	}

	if errCount > allowableErrorCount {
		return ErrURLContainsInvalidHost
	}

	return nil
}

func (v validation) ValidateRepoFiles(path string) error {
	var errCount int

	files, err := getRepositoryFiles(path)
	if err != nil {
		if errors.Is(err, ErrDirectoryNotFound) {
			return err
		}

		return fmt.Errorf("get list of repository files: %w", err)
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

	_, err := os.Lstat(path)
	if err != nil {
		return "", ErrDirectoryNotFound
	}

	err = filepath.Walk(path, func(_ string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk through files in directory: %w", err)
		}

		files += fmt.Sprintf(" %s", info.Name())

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("get list of repository files: %w", err)
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

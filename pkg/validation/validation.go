package validation

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

type validation struct{}

var _ Validator = (*validation)(nil)

func New() *validation {
	return &validation{}
}

func (v validation) ValidateURL(url string) error {
	hosts := [2]string{"github.com", "gitlab.com"}
	errCount := 0

	for _, host := range hosts {
		if strings.Contains(url, host) {
			continue
		}

		errCount++
	}

	if errCount > 1 {
		return fmt.Errorf("url contains invalid host: %s", url)
	}

	return nil
}

func (v validation) ValidateRepoFiles(path string) error {
	var data string

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to get repository files: %w", err)
		}

		data += fmt.Sprintf(" %s", info.Name())

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to get repository files: %w", err)
	}

	baseFiles := [2]string{"init.lua", "init.vim"}
	errCount := 0
	for _, file := range baseFiles {
		if strings.Contains(data, file) {
			continue
		}
		errCount++
	}

	if errCount > 1 {
		return fmt.Errorf("repository didn't contains base files for nvim configuration")
	}

	return nil
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

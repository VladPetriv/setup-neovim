package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

func (s service) CloneAndValidateRepository(url string) error {
	if isContainsSshURL(url) {
		return fmt.Errorf("sorry we don't support ssh urls")
	}

	_, err := git.PlainClone("nvim", false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	err = s.validator.ValidateRepoFiles("nvim")
	if err != nil {
		removeErr := os.RemoveAll("nvim")
		if removeErr != nil {
			return fmt.Errorf("repository validation failed, failed to remove repository: %w", err)
		}

		return fmt.Errorf("repository validation failed: %w", err)
	}

	return nil
}

func isContainsSshURL(url string) bool {
	return strings.Contains(url, "git@")
}

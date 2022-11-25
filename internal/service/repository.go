package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func (s service) CloneAndValidateRepository(url string) error {
	options := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	if isContainsSshURL(url) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		keyPath, err := s.input.GetInput("Enter path to your ssh file(.ssh/id_ed3122)")
		if err != nil {
			return fmt.Errorf("failed to get user input: %w", err)
		}

		filePath := fmt.Sprintf("%s/%s", homeDir, keyPath)

		publicKeys, err := ssh.NewPublicKeysFromFile("git", filePath, "")
		if err != nil {
			return fmt.Errorf("failed to create public keys from file: %w", err)
		}

		options.Auth = publicKeys
	}

	_, err := git.PlainClone("nvim", false, options)
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

package service

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/go-git/go-git/v5"
	gogitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type repositoryService struct {
	validator validation.Validator
}

func NewRepository(validator validation.Validator) RepositoryService {
	return &repositoryService{validator: validator}
}

// HasSSHURL reports whether the given URL uses SSH (git@) transport.
func HasSSHURL(url string) bool {
	return strings.Contains(url, "git@")
}

func (r repositoryService) CloneRepository(url string, sshKeyPath string) error {
	cloneOptions := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	if HasSSHURL(url) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home directory: %w", err)
		}

		publicKeys, err := gogitssh.NewPublicKeysFromFile("git", fmt.Sprintf("%s/%s", homeDir, sshKeyPath), "")
		if err != nil {
			return fmt.Errorf("create ssh public keys: %w", err)
		}

		cloneOptions.Auth = publicKeys
	}

	_, err := git.PlainClone(DirectoryNameForClonnedRepository, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}

	return nil
}

func (r repositoryService) ValidateRepository(path string) error {
	err := r.validator.ValidateRepoFiles(path)
	if err != nil {
		removeErr := os.RemoveAll(path)
		if removeErr != nil {
			return fmt.Errorf("repository validation failed, failed to remove repository: %w", err)
		}

		if errors.Is(err, validation.ErrNoBaseFilesInRepository) {
			return err
		}

		return fmt.Errorf("validate repository files: %w", err)
	}

	return nil
}

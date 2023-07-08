package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func (s service) CloneAndValidateRepository(url string, stdin io.Reader) error {
	cloneOptions := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	if hasSSHURLParts(url) {
		publicKeys, err := createPublicSSHKeysFromFile(s.input, stdin)
		if err != nil {
			return fmt.Errorf("failed to process ssh url: %w", err)
		}

		cloneOptions.Auth = publicKeys
	}

	_, err := git.PlainClone("nvim", false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	err = s.validator.ValidateRepoFiles("nvim")
	if err != nil {
		removeErr := os.RemoveAll("nvim")
		if removeErr != nil {
			return fmt.Errorf("repository validation failed, failed to remove repository: %w", err)
		}

		if errors.Is(err, validation.ErrNoBaseFilesInRepository) {
			return err
		}

		return fmt.Errorf("repository validation failed: %w", err)
	}

	return nil
}

func hasSSHURLParts(url string) bool {
	return strings.Contains(url, "git@")
}

func createPublicSSHKeysFromFile(input input.Inputter, stdin io.Reader) (*ssh.PublicKeys, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	fmt.Print("Enter path to your ssh file(.ssh/id_ed3122): ")

	keyPath, err := input.GetInput(stdin)
	if err != nil {
		return nil, fmt.Errorf("get input for ssh key path: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s", homeDir, keyPath)

	publicKeys, err := ssh.NewPublicKeysFromFile("git", filePath, "")
	if err != nil {
		return nil, fmt.Errorf("create ssh public from file: %w", err)
	}

	return publicKeys, nil
}

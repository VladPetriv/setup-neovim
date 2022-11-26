package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func (s service) CloneAndValidateRepository(url string) error {
	cloneOptions := &git.CloneOptions{ //nolint
		URL:      url,
		Progress: os.Stdout,
	}

	if haveSSHURLParts(url) {
		publicKeys, err := createPublicSSHKeysFromFile(s.input)
		if err != nil {
			return fmt.Errorf("failed process ssh url: %w", err)
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

		return fmt.Errorf("repository validation failed: %w", err)
	}

	return nil
}

func haveSSHURLParts(url string) bool {
	return strings.Contains(url, "git@")
}

func createPublicSSHKeysFromFile(input input.Inputter) (*ssh.PublicKeys, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	fmt.Print("Enter path to your ssh file(.ssh/id_ed3122): ") //nolint

	keyPath, err := input.GetInput(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to get user input: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s", homeDir, keyPath)

	publicKeys, err := ssh.NewPublicKeysFromFile("git", filePath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create ssh public keys from file: %w", err)
	}

	return publicKeys, nil
}

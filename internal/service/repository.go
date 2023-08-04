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

type repositoryService struct {
	inputter  input.Inputter
	validator validation.Validator
}

func NewRepository(inputter input.Inputter, validator validation.Validator) *repositoryService {
	return &repositoryService{
		inputter:  inputter,
		validator: validator,
	}
}

func (r repositoryService) ProcessUserURL(stdin io.Reader) (string, error) {
	fmt.Print("Enter URL to your nvim config: ")

	configURL, err := r.inputter.GetInput(stdin)
	if err != nil {
		return "", fmt.Errorf("get user input: %w", err)
	}

	err = r.validator.ValidateURL(configURL)
	if err != nil {
		return "", err
	}

	return configURL, nil
}

func (r repositoryService) CloneAndValidateRepository(url string, stdin io.Reader) error {
	cloneOptions := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	if hasSSHURLParts(url) {
		publicKeys, err := createPublicSSHKeysFromFile(r.inputter, stdin)
		if err != nil {
			return fmt.Errorf("create public ssh key from file: %w", err)
		}

		cloneOptions.Auth = publicKeys
	}

	_, err := git.PlainClone("nvim", false, cloneOptions)
	if err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}

	err = r.validator.ValidateRepoFiles("nvim")
	if err != nil {
		removeErr := os.RemoveAll("nvim")
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
		return nil, fmt.Errorf("get user input: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s", homeDir, keyPath)

	publicKeys, err := ssh.NewPublicKeysFromFile("git", filePath, "")
	if err != nil {
		return nil, fmt.Errorf("create ssh public keys: %w", err)
	}

	return publicKeys, nil
}

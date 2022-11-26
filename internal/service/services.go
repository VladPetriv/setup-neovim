package service

import (
	"io"

	"github.com/VladPetriv/setup-neovim/internal/models"
)

type Services interface {
	// CheckUtilStatus check if nvim and git are installed
	CheckUtilStatus() map[string]string
	// ProcessUserURL get URL by user input and validate them
	ProcessUserURL(stdin io.Reader) (string, error)
	// CloneAndValidateRepository clones git repository and check if repository have base files for nvim configuration
	CloneAndValidateRepository(url string, stdin io.Reader) error
	// MoveConfigDirectory moves repository directory into .config directory
	MoveConfigDirectory() error
	// ProcessPackageManagers ask user about package managers and install them if needed
	ProcessPackageManagers(stdin io.Reader) (models.PackageManager, error)
	// CompleteSetup runs nvim with specific flag that depends on package manager
	CompleteSetup(packageManager models.PackageManager) error
}

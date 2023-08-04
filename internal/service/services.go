package service

import (
	"errors"
	"io"
)

type Services interface {
	// CheckUtilStatus check if nvim and git are installed
	CheckUtilStatus() map[string]string
	// ProcessUserURL get URL by user input and validate them
	ProcessUserURL(stdin io.Reader) (string, error)
	// CloneAndValidateRepository clones git repository and check if repository have base files for nvim configuration
	CloneAndValidateRepository(url string, stdin io.Reader) error
	// InstallPackageManager ask user about which package manager to install
	InstallPackageManager(stdin io.Reader) (string, error)
	// ExtractAndMoveConfigDirectory get config directory from repository and move it to .config folder
	ExtractAndMoveConfigDirectory(path string) error
	// DeleteConfigOrStopInstallation checks if nvim config is exist and ask permission for deleting it.
	DeleteConfigOrStopInstallation(stdin io.Reader) error
	// DetectInstalledPackageManagers check if user has installed of any package managers.
	DetectInstalledPackageManagers() (string, int, error)
	// ProcessAlreadyInstalledPackageManagers is used to inform user about already installed managers
	// and ask permission for deleting old and installing new.
	ProcessAlreadyInstalledPackageManagers(countOfAlreadyInstalledManagers int, stdin io.Reader) (bool, error)
}

type ManagerService interface {
	// DetectInstalledPackageManagers check if user has already installed package managers.
	DetectInstalledPackageManagers() (string, int, error)
	// ProcessAlreadyInstalledPackageManagers inform user about already installed package managers
	// and ask user permission for deleting them and deletet them if user want it.
	ProcessAlreadyInstalledPackageManagers(countOfAlreadyInstalledManagers int, stdin io.Reader) (bool, error)
	// InstallPackageManager ask user about which package manager to install and install it.
	InstallPackageManager(stdin io.Reader) (string, error)
}

type RepositoryService interface {
	// ProcessUserURL get URL by user input and validate them
	ProcessUserURL(stdin io.Reader) (string, error)
	// CloneAndValidateRepository clones git repository and check if repository have base files for nvim configuration
	CloneAndValidateRepository(url string, stdin io.Reader) error
}

var (
	ErrDirectoryNotFound     = errors.New("directory not found")
	ErrEnterValidAnswer      = errors.New("please enter valid answer")
	ErrDirectoryAlreadyExist = errors.New("config directory already exists")
	ErrStopInstallation      = errors.New("stop config installation")
	ErrConfigNotFound        = errors.New("nvim config not found")
	ErrNoNeedToDelete        = errors.New("not need to delete")
)

// TODO: Create a custom type for package managers

const (
	PackerPluginManager  = "packer"
	VimPlugPluginManager = "vim-plug"
)

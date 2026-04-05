package service

import "errors"

type Services struct {
	File       FileService
	Manager    ManagerService
	Repository RepositoryService
}

type FileService interface {
	// CheckUtilStatus checks if nvim and git are installed.
	CheckUtilStatus() map[string]string
	// CheckConfigExists reports whether ~/.config/nvim already exists.
	CheckConfigExists() (bool, error)
	// DeleteConfig removes the existing ~/.config/nvim directory.
	DeleteConfig() error
	// ExtractAndMoveConfigDirectory gets the config directory from the cloned repo and moves it to ~/.config/nvim.
	ExtractAndMoveConfigDirectory(path string) error
	// DeleteRepositoryDirectory removes the cloned repository directory.
	DeleteRepositoryDirectory(path string) error
}

type RepositoryService interface {
	// CloneRepository clones the git repository at url. sshKeyPath is only used for SSH URLs.
	CloneRepository(url string, sshKeyPath string) error
	// ValidateRepository checks that the cloned repo contains init.lua or init.vim.
	ValidateRepository(path string) error
}

type ManagerService interface {
	// DetectInstalledPackageManagers returns the names of any already-installed package managers.
	DetectInstalledPackageManagers() ([]string, error)
	// DeletePackageManagers removes all detected package manager directories.
	DeletePackageManagers() error
	// InstallPackageManager installs the package manager identified by name.
	InstallPackageManager(name string) error
}

var ErrDirectoryNotFound = errors.New("directory not found")

const (
	PackerPluginManager  = "packer"
	VimPlugPluginManager = "vim-plug"
	NonePluginManager    = "None"

	DirectoryNameForClonnedRepository = "nvim"
	systemNeovimConfigPath            = ".config/nvim"
)

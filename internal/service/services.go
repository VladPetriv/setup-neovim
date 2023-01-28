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
	// ProcessPackageManagers ask user about package managers and install them if needed
	ProcessPackageManagers(stdin io.Reader) (string, error)
	// ExtractAndMoveConfigDirectory get config directory from repository and move it to .config folder
	ExtractAndMoveConfigDirectory(path string) error
}

var (
	ErrDirectoryNotFound = errors.New("directory not found")
	ErrEnterValidAnswer  = errors.New("please enter valid answer")
)

const (
	PackerPluginManager  = "packer"
	VimPlugPluginManager = "vim-plug"
)

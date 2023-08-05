package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/VladPetriv/setup-neovim/pkg/input"
)

type managerService struct {
	inputter input.Inputter
}

var _ ManagerService = (*managerService)(nil)

func NewManager(inputter input.Inputter) ManagerService {
	return &managerService{
		inputter: inputter,
	}
}

func (m managerService) DetectInstalledPackageManagers() (string, int, error) {
	result := map[string]bool{
		PackerPluginManager:  false,
		VimPlugPluginManager: false,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", 0, fmt.Errorf("get home directorY: %w", err)
	}

	if _, err = os.Lstat(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir)); err == nil {
		result[VimPlugPluginManager] = true
	}

	if _, err = os.Lstat(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir)); err == nil {
		result[PackerPluginManager] = true
	}

	var count int
	var message string

	for manager, installed := range result {
		if !installed {
			continue
		}

		count++
		message += fmt.Sprintf("Detected already installed package manager!\nPackage manager name: %s\n", manager)
	}

	return message, count, nil
}

func (m managerService) ProcessAlreadyInstalledPackageManagers(count int, stdin io.Reader) (bool, error) {
	if count == 0 {
		return true, nil
	}

	err := m.deletePackageManagers(stdin)
	if err != nil {
		if errors.Is(err, ErrNoNeedToDelete) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// TODO: This function should not be as structure method.
func (m managerService) deletePackageManagers(stdin io.Reader) error {
	fmt.Print("Do you want to remove old package managers and install new?(y/n): ")

	reader := bufio.NewReader(stdin)
	wantRemove, err := m.inputter.GetInput(reader)
	if err != nil {
		return fmt.Errorf("get user input: %w", err)
	}

	// TODO: We need to accept a map with already installed to delete only installed managers
	// instead of trying to delete both
	switch wantRemove {
	case "y":
		homeDir, homeDirErr := os.UserHomeDir()
		if homeDirErr != nil {
			return fmt.Errorf("get home directory: %w", err)
		}

		if err = os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir)); err != nil {
			return fmt.Errorf("delete vim-plug: %w", err)
		}

		if err = os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir)); err != nil {
			return fmt.Errorf("delete packer: %w", err)
		}

		return nil

	case "n":
		return ErrNoNeedToDelete
	default:
		return ErrEnterValidAnswer
	}
}

func (m managerService) InstallPackageManager(stdin io.Reader) (string, error) {
	fmt.Print("Choose package manager(packer/vim-plug): ")

	reader := bufio.NewReader(stdin)
	packageManager, err := m.inputter.GetInput(reader)
	if err != nil {
		return "", fmt.Errorf("get user input: %w", err)
	}

	switch packageManager {
	case PackerPluginManager:
		err = installPacker()
		if err != nil {
			return "", fmt.Errorf("install packer: %w", err)
		}

		return PackerPluginManager, nil
	case VimPlugPluginManager:
		err = installVimPlug()
		if err != nil {
			return "", fmt.Errorf("install vim-plug: %w", err)
		}

		return VimPlugPluginManager, nil
	default:
		return "", ErrEnterValidAnswer
	}
}

func installVimPlug() error {
	file, err := os.Create("vim-plug.sh")
	if err != nil {
		return fmt.Errorf("create file for installing script: %w", err)
	}

	_, err = file.Write([]byte(
		`
      #!/bin/sh
      curl -fLo "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/autoload/plug.vim --create-dirs \
        https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
    `,
	))
	if err != nil {
		return fmt.Errorf("write data to installing script file: %w", err)
	}

	cmd := exec.Command("/bin/sh", "vim-plug.sh")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("execute installation script file: %w", err)
	}

	err = os.Remove("./vim-plug.sh")
	if err != nil {
		return fmt.Errorf("remove installation script file: %w", err)
	}

	return nil
}

func installPacker() error {
	file, err := os.Create("packer.sh")
	if err != nil {
		return fmt.Errorf("create file for installing script: %w", err)
	}

	_, err = file.Write([]byte(
		`
      #!/bin/sh
      git clone --depth 1 https://github.com/wbthomason/packer.nvim\
        ~/.local/share/nvim/site/pack/packer/start/packer.nvim
    `,
	))
	if err != nil {
		return fmt.Errorf("write data to installing script file: %w", err)
	}

	cmd := exec.Command("/bin/sh", "packer.sh")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("execute installation script file: %w", err)
	}

	err = os.Remove("./packer.sh")
	if err != nil {
		return fmt.Errorf("remove installation script file: %w", err)
	}

	return nil
}

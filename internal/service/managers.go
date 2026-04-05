package service

import (
	"fmt"
	"os"
	"os/exec"
)

type managerService struct{}

func NewManager() ManagerService {
	return &managerService{}
}

func (m managerService) DetectInstalledPackageManagers() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	var detected []string

	if _, err = os.Lstat(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir)); err == nil {
		detected = append(detected, VimPlugPluginManager)
	}

	if _, err = os.Lstat(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir)); err == nil {
		detected = append(detected, PackerPluginManager)
	}

	return detected, nil
}

func (m managerService) DeletePackageManagers() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	if err = os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir)); err != nil {
		return fmt.Errorf("delete vim-plug: %w", err)
	}

	if err = os.RemoveAll(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir)); err != nil {
		return fmt.Errorf("delete packer: %w", err)
	}

	return nil
}

func (m managerService) InstallPackageManager(name string) error {
	switch name {
	case PackerPluginManager:
		return installPacker()
	case VimPlugPluginManager:
		return installVimPlug()
	case NonePluginManager:
		return nil
	default:
		return fmt.Errorf("unknown package manager: %s", name)
	}
}

func installVimPlug() error {
	file, err := os.Create("vim-plug.sh")
	if err != nil {
		return fmt.Errorf("create installation script file: %w", err)
	}

	_, err = file.Write([]byte(
		`
      #!/bin/sh
      curl -fLo "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/autoload/plug.vim --create-dirs \
        https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
    `,
	))
	if err != nil {
		return fmt.Errorf("write data to installation script file: %w", err)
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
		return fmt.Errorf("create installation script file: %w", err)
	}

	_, err = file.Write([]byte(
		`
      #!/bin/sh
      git clone --depth 1 https://github.com/wbthomason/packer.nvim\
        ~/.local/share/nvim/site/pack/packer/start/packer.nvim
    `,
	))
	if err != nil {
		return fmt.Errorf("write data to installation script file: %w", err)
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

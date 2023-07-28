package service

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (s service) InstallPackageManager(stdin io.Reader) (string, error) {
	fmt.Print("Choose package manager(packer/vim-plug): ")

	reader := bufio.NewReader(stdin)
	packageManager, err := s.input.GetInput(reader)
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

func (s service) GetPackageMangerIfNotInstalled(stdin io.Reader) (string, error) {
	fmt.Print("Do you have any package managers installed?(y/n): ")

	reader := bufio.NewReader(stdin)
	haveInstalledPackageManagers, err := s.input.GetInput(reader)
	if err != nil {
		return "", fmt.Errorf("get input for is user has installed pkg manager: %w", err)
	}

	switch haveInstalledPackageManagers {
	case "y":
		return "skip", nil
	case "n":
		fmt.Print("Choose package manager(packer/vim-plug): ")

		return s.input.GetInput(reader)
	default:
		return "", nil
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
		return fmt.Errorf("execute script for vim plug installation: %w", err)
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
		return fmt.Errorf("execute script for packer installation: %w", err)
	}

	err = os.Remove("./packer.sh")
	if err != nil {
		return fmt.Errorf("remove installation script file: %w", err)
	}

	return nil
}

func (s service) DetectInstalledPackageManagers() (map[string]bool, error) {
	result := map[string]bool{
		PackerPluginManager:  false,
		VimPlugPluginManager: false,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directorY: %w", err)
	}

	if _, err = os.Lstat(fmt.Sprintf("%s/.local/share/nvim/site/autoload", homeDir)); err == nil {
		result[VimPlugPluginManager] = true
	}

	if _, err = os.Lstat(fmt.Sprintf("%s/.local/share/nvim/site/pack", homeDir)); err == nil {
		result[PackerPluginManager] = true
	}

	return result, nil
}

func (s service) ProcessAlreadyInstalledPackageManagers(stdin io.Reader) {}

func (s service) DeletePackageManagers(stdin io.Reader) error {
	fmt.Print("Do you want to remove old package managers and install new?(y/n): ")

	reader := bufio.NewReader(stdin)
	wantRemove, err := s.input.GetInput(reader)
	if err != nil {
		return fmt.Errorf("get user input: %w", err)
	}

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

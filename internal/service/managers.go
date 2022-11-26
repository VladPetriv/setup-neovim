package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/VladPetriv/setup-neovim/internal/models"
	"github.com/VladPetriv/setup-neovim/pkg/input"
)

var ErrEnterValidAnswer = errors.New("please enter valid answer")

func (s service) ProcessPackageManagers() (models.PackageManager, error) {
	fmt.Print("Do you have any package managers installed?(y/n): ") //nolint

	haveInstalledPackageManagers, err := s.input.GetInput()
	if err != nil {
		return "", fmt.Errorf("failed to get user input: %w", err)
	}

	switch haveInstalledPackageManagers {
	case "n":
		return installPackageManager(s.input)
	case "y":
		return "", nil
	default:
		return "", ErrEnterValidAnswer
	}
}

func installPackageManager(input input.Inputter) (models.PackageManager, error) {
	fmt.Print("Choose package manager(packer/vim-plug): ") //nolint

	packageManager, err := input.GetInput()
	if err != nil {
		return "", fmt.Errorf("failed to get user input: %w", err)
	}

	switch models.PackageManager(packageManager) {
	case models.Packer:
		err := installPacker()
		if err != nil {
			return "", fmt.Errorf("install packer error: %w", err)
		}

		return models.Packer, nil
	case models.VimPlug:
		err := installVimPlug()
		if err != nil {
			return "", fmt.Errorf("install vim-plug error: %w", err)
		}

		return models.VimPlug, nil
	default:
		return "", ErrEnterValidAnswer
	}
}

func installVimPlug() error {
	file, err := os.Create("vim-plug.sh")
	if err != nil {
		return fmt.Errorf("create sh file for install script error: %w", err)
	}

	_, err = file.Write([]byte(
		`
      #!/bin/sh
      curl -fLo "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/autoload/plug.vim --create-dirs \
        https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
    `,
	))
	if err != nil {
		return fmt.Errorf("write to sh file error: %w", err)
	}

	cmd := exec.Command("/bin/sh", "vim-plug.sh")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("run sh file error: %w", err)
	}

	err = os.Remove("./vim-plug.sh")
	if err != nil {
		return fmt.Errorf("remove sh file error: %w", err)
	}

	return nil
}

func installPacker() error {
	file, err := os.Create("packer.sh")
	if err != nil {
		return fmt.Errorf("create file for install script error: %w", err)
	}

	_, err = file.Write([]byte(
		`
      #!/bin/sh
      git clone --depth 1 https://github.com/wbthomason/packer.nvim\
        ~/.local/share/nvim/site/pack/packer/start/packer.nvim
    `,
	))
	if err != nil {
		return fmt.Errorf("write to sh file error: %w", err)
	}

	cmd := exec.Command("/bin/sh", "packer.sh")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("run sh file error: %w", err)
	}

	err = os.Remove("./packer.sh")
	if err != nil {
		return fmt.Errorf("remove sh file error: %w", err)
	}

	return nil
}

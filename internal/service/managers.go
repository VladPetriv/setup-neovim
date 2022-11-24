package service

import (
	"fmt"
	"os"
	"os/exec"
)

func (s service) ProcessPackageManagers() (string, error) {
	haveInstalledPackageManagers, err := s.input.GetInput("Do you have any package managers installed?(y/n)")
	if err != nil {
		return "", fmt.Errorf("failed to get user input: %w", err)
	}

	switch haveInstalledPackageManagers {
	case "n":
		packageManager, err := s.input.GetInput("Choose package manager(packer/vim-plug)")
		if err != nil {
			return "", fmt.Errorf("failed to get user input: %w", err)
		}

		switch packageManager {
		case "packer":
			err := installPacker()
			if err != nil {
				return "", fmt.Errorf("install packer: %w", err)
			}
			return "packer", nil

		case "vim-plug":
			err := installVimPlug()
			if err != nil {
				return "plug", fmt.Errorf("install vim-plug: %w", err)
			}
			return "", nil
		default:
			return "", fmt.Errorf("please choose valid package manager(packer/vim-plug): %s", packageManager)
		}
	case "y":
		return "", nil
	default:
		return "", fmt.Errorf("please enter valid answer(y/n): %s", haveInstalledPackageManagers)
	}
}

func installVimPlug() error {
	file, err := os.Create("vim-plug.sh")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write([]byte(
		`
      #!/bin/sh
      curl -fLo "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/autoload/plug.vim --create-dirs \
        https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
    `,
	))
	if err != nil {
		return fmt.Errorf("failed to create sh file for installing vim-plug: %w", err)
	}

	cmd := exec.Command("/bin/sh", "vim-plug.sh")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to install vim-plug: %w", err)
	}

	err = os.Remove("./vim-plug.sh")
	if err != nil {
		return fmt.Errorf("failed to remove file after installing: %w", err)
	}

	return nil
}

func installPacker() error {
	file, err := os.Create("packer.sh")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write([]byte(
		`
      #!/bin/sh
      git clone --depth 1 https://github.com/wbthomason/packer.nvim\
        ~/.local/share/nvim/site/pack/packer/start/packer.nvim
    `,
	))
	if err != nil {
		return fmt.Errorf("failed to create sh file for installing packer: %w", err)
	}

	cmd := exec.Command("/bin/sh", "packer.sh")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to install packer: %w", err)
	}

	err = os.Remove("./packer.sh")
	if err != nil {
		return fmt.Errorf("failed to remove file after installing: %w", err)
	}

	return nil
}

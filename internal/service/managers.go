package service

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (s service) ProcessPackageManagers(stdin io.Reader) (string, error) {
	packageManager, err := s.input.ProcessInputForPackageManagers(stdin)
	if err != nil {
		return "", fmt.Errorf("failed to process user input: %w", err)
	}

	switch packageManager {
	case PackerPluginManager:
		err = installPacker()
		if err != nil {
			return "", fmt.Errorf("install packer error: %w", err)
		}

		return PackerPluginManager, nil
	case VimPlugPluginManager:
		err = installVimPlug()
		if err != nil {
			return "", fmt.Errorf("install vim-plug error: %w", err)
		}
		return VimPlugPluginManager, nil
	case "skip":
		return "", nil
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

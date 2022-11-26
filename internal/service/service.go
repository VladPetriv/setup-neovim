package service

import (
	"fmt"
	"os/exec"

	"github.com/VladPetriv/setup-neovim/internal/models"
	"github.com/VladPetriv/setup-neovim/pkg/input"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

type service struct {
	validator validation.Validator
	input     input.Inputter
}

var _ Services = (*service)(nil)

type Options struct {
	Validator validation.Validator
	Inputter  input.Inputter
}

func New(options *Options) *service { //nolint
	return &service{
		validator: options.Validator,
		input:     options.Inputter,
	}
}

func (s service) CheckUtilStatus() map[string]string {
	return s.validator.ValidateConsoleUtilities()
}

func (s service) CompleteSetup(packageManager models.PackageManager) error {
	var args string

	if packageManager == models.Packer {
		args = "+PackerSync"
	}

	if packageManager == models.VimPlug {
		args = "+PlugInstall"
	}

	cmd := exec.Command("nvim", args)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to complete setup: %w", err)
	}

	return nil
}

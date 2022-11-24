package service

import (
	"fmt"
	"os"
)

func (s service) MoveConfigDirectory() error {
	homeeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get path to the home directory: %w", err)
	}

	err = os.Rename("./nvim", fmt.Sprintf("%s/.config/nvim", homeeDir))
	if err != nil {
		return fmt.Errorf("failed to move repository into .config directory: %w", err)
	}

	return nil
}
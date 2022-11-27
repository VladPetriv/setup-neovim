package service

import (
	"fmt"
	"os"
)

func (s service) MoveConfigDirectory() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get path to the home directory: %w", err)
	}

	err = os.Rename("./nvim", fmt.Sprintf("%s/.config/nvim", homeDir))
	if err != nil {
		removeErr := os.RemoveAll("./nvim")
		if removeErr != nil {
			return fmt.Errorf("moving config failed, failed to remove repository: %w", err)
		}

		return fmt.Errorf("failed to move repository into .config directory: %w", err)
	}

	return nil
}

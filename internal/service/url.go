package service

import (
	"fmt"
	"os"
)

func (s service) ProcessUserURL() (string, error) {
	fmt.Print("Enter URL to your nvim config: ") //nolint

	configURL, err := s.input.GetInput(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed get user input: %w", err)
	}

	err = s.validator.ValidateURL(configURL)
	if err != nil {
		return "", fmt.Errorf("URL validation failed: %w", err)
	}

	return configURL, nil
}

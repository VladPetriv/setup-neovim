package service

import (
	"fmt"
	"io"
)

func (s service) ProcessUserURL(stdin io.Reader) (string, error) {
	fmt.Print("Enter URL to your nvim config: ")

	configURL, err := s.input.GetInput(stdin)
	if err != nil {
		return "", fmt.Errorf("failed get user input: %w", err)
	}

	err = s.validator.ValidateURL(configURL)
	if err != nil {
		return "", fmt.Errorf("URL validation failed: %w", err)
	}

	return configURL, nil
}

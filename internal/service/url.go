package service

import "fmt"

func (s service) ProcessUserURL() (string, error) {
	configURL, err := s.input.GetInput("Enter URL to your nvim config")
	if err != nil {
		return "", fmt.Errorf("failed get user input: %w", err)
	}

	err = s.validator.ValidateURL(configURL)
	if err != nil {
		return "", fmt.Errorf("validation for URL failed: %w", err)
	}

	return configURL, nil
}

package app

import (
	"fmt"
	"os"
	"time"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/colors"
)

func Run(service service.Services) {
	// commandTimeout represents timeout that should be after completing the previous function.
	commandTimeout := 1 * time.Second

	errs := service.CheckUtilStatus()
	if len(errs) >= 1 {
		colors.Red(fmt.Sprintf("Errors: %v\n", errs))
		os.Exit(1)
	}

	colors.Green("All utilities are installed....")

	time.Sleep(commandTimeout)

	url, err := service.ProcessUserURL(os.Stdin)
	if err != nil {
		colors.Red("Validation for URL failed! Please try again... ")
		colors.Red(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	colors.Green("URL is valid...")

	time.Sleep(commandTimeout)

	err = service.CloneAndValidateRepository(url, os.Stdin)
	if err != nil {
		colors.Red("Failed to clone repository or repository didn't have base files for nvim configuration")
		colors.Red(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	colors.Green("Repository successfully cloned and checked for base files")

	time.Sleep(commandTimeout)

	err = service.ExtractAndMoveConfigDirectory("./nvim")
	if err != nil {
		colors.Red("Failed to extract and move config directory from repository")
		colors.Red(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}
	colors.Green("Config successfully extracted!")

	time.Sleep(commandTimeout)

	packageManger, err := service.ProcessPackageManagers(os.Stdin)
	if err != nil {
		colors.Red("Failed to install package managers. Please try again...")
		colors.Red(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	if packageManger != "" {
		colors.Green(fmt.Sprintf("%s successfully installed", packageManger))

		time.Sleep(commandTimeout)
	}

	colors.Green("Nvim successfully configured!")
}

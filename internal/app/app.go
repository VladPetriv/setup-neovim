package app

import (
	"fmt"
	"os"
	"time"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/colors"
)

// commandTimeout represents timeout that should be after completing the previous function.
var commandTimeout = 1 * time.Second //nolint

func Run(service service.Services) { //nolint
	errs := service.CheckUtilStatus()
	if len(errs) >= 1 {
		colors.Red(fmt.Sprintf("Errors: %v\n", errs))
		os.Exit(1)
	}

	colors.Green("All utilities are installed....")

	time.Sleep(commandTimeout)

	url, err := service.ProcessUserURL()
	if err != nil {
		colors.Red("Validation for URL failed! Please try again... ")
		colors.Red(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	colors.Green("URL is valid...")

	time.Sleep(commandTimeout)

	err = service.CloneAndValidateRepository(url)
	if err != nil {
		colors.Red("Failed to clone repository or repository didn't have base files for nvim configuration")
		colors.Red(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	colors.Green("Repository successfully cloned and checked for base files")

	time.Sleep(commandTimeout)

	err = service.MoveConfigDirectory()
	if err != nil {
		colors.Red("Failed to move repository to .config directory! Please try again... ")
		colors.Red(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	colors.Green("Successfully moved repository to .config directory...")

	time.Sleep(commandTimeout)

	packageManger, err := service.ProcessPackageManagers()
	if err != nil {
		colors.Red("Failed to install package managers. Please try again...")
		colors.Red(fmt.Sprintf("Error: %v\n", err))
		os.Exit(1)
	}

	colors.Green(fmt.Sprintf("%s successfully installed", packageManger))

	time.Sleep(commandTimeout)

	err = service.CompleteSetup(packageManger)
	if err != nil {
		colors.Red(
			fmt.Sprintf("Failed to run nvim with %s installation command! Please try to run it manually... ",
				packageManger,
			),
		)
		os.Exit(1)
	}

	colors.Green("Nvim successfully configured!")
}

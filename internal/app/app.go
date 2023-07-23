package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/VladPetriv/setup-neovim/internal/service"
	"github.com/VladPetriv/setup-neovim/pkg/colors"
	"github.com/VladPetriv/setup-neovim/pkg/errs"
	"github.com/VladPetriv/setup-neovim/pkg/validation"
)

// commandTimeout represents timeout that should be after completing the previous function.
const commandTimeout = 1 * time.Second

func Run(srv service.Services) {
	utilErrors := srv.CheckUtilStatus()
	if len(utilErrors) >= 1 {
		for _, value := range utilErrors {
			colors.Red(value)
		}

		os.Exit(1)
	}

	successfulCommand("All utilities are installed....")

	err := srv.DeleteConfigOrStopInstallation(os.Stdin)
	if err != nil && !errors.Is(err, service.ErrConfigNotFound) {
		if errors.Is(err, service.ErrStopInstallation) {
			errs.WrapError("Thank you for using setup-nvim!", err)
		}

		if errors.Is(err, service.ErrEnterValidAnswer) {
			errs.WrapError("Please choose correct answer for question!", err)
		}

		errs.WrapError("Failed to delete config or stop installation!", err)
	}
	if !errors.Is(err, service.ErrConfigNotFound) {
		successfulCommand("Successfully remove old nvim config...")
	}

	url, err := srv.ProcessUserURL(os.Stdin)
	if err != nil {
		errs.WrapError("Validation for URL failed! Please try again... ", err)
	}

	successfulCommand("URL is valid...")

	err = srv.CloneAndValidateRepository(url, os.Stdin)
	if err != nil {
		if errors.Is(err, validation.ErrNoBaseFilesInRepository) {
			errs.WrapError("Repository didn't have base files for nvim configuration!", err)
		}

		errs.WrapError("Failed to clone repository or validate repository!", err)
	}

	successfulCommand("Repository successfully cloned and checked for base files")

	err = srv.ExtractAndMoveConfigDirectory("./nvim")
	if err != nil {
		errs.WrapError("Failed to extract and move config directory from repository", err)
	}

	successfulCommand("Config successfully extracted and moved to '.config' directory!")

	alreadyInstalledManagers, err := srv.DetectInstalledPackageManagers()
	if err != nil {
		errs.WrapError("Failed to detect installed package managers", err)
	}

	var alreadyInstalledManagersCount int

	for manager, installed := range alreadyInstalledManagers {
		if installed {
			alreadyInstalledManagersCount++
			colors.Yellow(
				fmt.Sprintf("Detected already installed package manager!\nPackage manager name: %s\n", manager),
			)
		}
	}

	if alreadyInstalledManagersCount > 0 {
		err = srv.DeletePackageManagersIfNeeded(os.Stdin)
		if err != nil {
			if errors.Is(err, service.ErrEnterValidAnswer) {
				errs.WrapError("Please choose correct answer for question!", err)
			}

			errs.WrapError("Failed to delete package managers!", err)
		}

		successfulCommand("Successfully remove all old package managers!")
	}

	packageManger, err := srv.InstallPackageManager(os.Stdin)
	if err != nil {
		errs.WrapError("Failed to install package manager. Please try again...", err)
	}

	if packageManger != "" {
		successfulCommand(fmt.Sprintf("%s successfully installed", packageManger))
	}

	colors.Green("Nvim successfully configured!")
}

// successfulCommand print message with green color and wait command timeout [3s].
func successfulCommand(text string) {
	colors.Green(text)

	time.Sleep(commandTimeout)
}

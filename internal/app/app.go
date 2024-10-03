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
	utilErrors := srv.File.CheckUtilStatus()
	if len(utilErrors) >= 1 {
		for _, value := range utilErrors {
			colors.Red(value)
		}

		os.Exit(0)
	}

	successfulMessage("All utilities are installed....")

	err := srv.File.DeleteConfigOrStopInstallation(os.Stdin)
	if err != nil && !errors.Is(err, service.ErrConfigNotFound) {
		if errors.Is(err, service.ErrStopInstallation) {
			errs.WrapError("Thank you for using setup-nvim!", err)
		}

		errs.WrapError("Failed to delete config or stop installation!", err)
	}
	if !errors.Is(err, service.ErrConfigNotFound) {
		successfulMessage("Successfully remove old nvim config...")
	}

	url, err := srv.Repository.ProcessUserURL(os.Stdin)
	if err != nil {
		errs.WrapError("Validation for URL failed! Please try again... ", err)
	}

	successfulMessage("URL is valid...")

	err = srv.Repository.CloneAndValidateRepository(url, os.Stdin)
	if err != nil {
		if errors.Is(err, validation.ErrNoBaseFilesInRepository) {
			errs.WrapError("Repository didn't have base files for nvim configuration!", err)
		}

		errs.WrapError("Failed to clone repository or validate repository!", err)
	}

	successfulMessage("Repository successfully cloned and checked for base files")

	repositoryPath := fmt.Sprintf("./%s", service.DirectoryNameForClonnedRepository)
	err = srv.File.ExtractAndMoveConfigDirectory(repositoryPath)
	if err != nil {
		deleteErr := srv.File.DeleteRepositoryDirectory(repositoryPath)
		if deleteErr != nil {
			warningMessage("Could not delete clonned repository")
		}

		errs.WrapError("Failed to extract and move config directory from repository", err)
	}

	deleteErr := srv.File.DeleteRepositoryDirectory(repositoryPath)
	if deleteErr != nil {
		warningMessage("Could not delete clonned repository")
	}

	successfulMessage("Config successfully extracted and moved to '.config' directory!")

	message, count, err := srv.Manager.DetectInstalledPackageManagers()
	if err != nil {
		errs.WrapError("Failed to detect installed package managers", err)
	}

	colors.Yellow(message)

	needToInstall, err := srv.Manager.ProcessAlreadyInstalledPackageManagers(count, os.Stdin)
	if err != nil {
		errs.WrapError("Failed to process already installed package managers!", err)
	}

	if !needToInstall {
		colors.Green("Nvim successfully configured!")

		os.Exit(0)
	}

	packageManger, installErr := srv.Manager.InstallPackageManager(os.Stdin)
	if installErr != nil {
		errs.WrapError("Failed to install package manager. Please try again...", err)
	}

	if packageManger != "" {
		successfulMessage(fmt.Sprintf("%s successfully installed", packageManger))
	}

	colors.Green("Nvim successfully configured!")
}

// successfulMessage print message with green color and wait input timeout.
func successfulMessage(text string) {
	colors.Green(text)

	time.Sleep(commandTimeout)
}

// warningMessage print message with yellow color.
func warningMessage(text string) {
	colors.Yellow(text)
}

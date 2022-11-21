package validation

type Validator interface {
	ValidateURL(url string) error
	ValidateRepoFiles(pathToRepo string) error
	ValidateConsoleUtilities() map[string]string
}

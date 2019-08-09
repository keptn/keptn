package common

import (
	"fmt"
	"strings"

	"github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/configuration-service/config"
)

// CheckoutBranch checks out the given branch
func CheckoutBranch(project string, branch string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := utils.ExecuteCommandInDirectory("git", []string{"checkout", branch}, projectConfigPath)
	if err != nil {
		return err
	}
	return nil
}

// StageAndCommitAll stages all current changes and commits them to the current branch
func StageAndCommitAll(project string, message string) error {
	projectConfigPath := config.ConfigDir + "/" + project
	_, err := utils.ExecuteCommandInDirectory("git", []string{"add", "."}, projectConfigPath)
	if err != nil {
		return err
	}

	out, err := utils.ExecuteCommandInDirectory("git", []string{"commit", "-m", `"` + message + `"`}, projectConfigPath)
	utils.Debug("", out)
	if err != nil && !(strings.Contains(err.Error(), "nothing to commit")) {
		fmt.Print(err.Error())
		return err
	}
	return nil
}

// GetCurrentVersion gets the latest version (i.e. commit hash) of the currently checked out branch
func GetCurrentVersion(project string) (string, error) {
	projectConfigPath := config.ConfigDir + "/" + project
	out, err := utils.ExecuteCommandInDirectory("git", []string{"rev-parse", "HEAD"}, projectConfigPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(out, "\n"), nil
}

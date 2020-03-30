package kube

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-version"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

var (
	kubeServerVersion  = regexp.MustCompile(`Server Version: version\.Info{Major:"(\d+)", Minor:"(\d+.*?)\+{0,1}"`)
	executeCommandFunc = keptnutils.ExecuteCommand
)

// IsKubectlAvailable checks whether kubectl is available
func IsKubectlAvailable() (bool, error) {

	_, err := executeCommandFunc("kubectl", []string{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func getKubeServerVersion() (string, error) {

	out, err := executeCommandFunc("kubectl", []string{"version"})
	if err != nil {
		return "", err
	}
	submatches := kubeServerVersion.FindStringSubmatch(out)
	if submatches == nil {
		return "", errors.New("Server Version not found: " + out)
	}
	return submatches[1] + "." + submatches[2], nil
}

// CheckKubeServerVersion checks the Kubernetes Server version against the given constraints
func CheckKubeServerVersion(constraints string) error {

	serverVersion, err := getKubeServerVersion()
	if err != nil {
		return err
	}
	hVersion, err := version.NewVersion(serverVersion)
	if err != nil {
		return err
	}
	hConstraints, err := version.NewConstraint(constraints)
	if err != nil {
		return err
	}
	if hConstraints.Check(hVersion) {
		return nil
	}
	return fmt.Errorf("The Kubernetes Server Version '%s' doesn't satisfy constraints '%s'", serverVersion, constraints)
}

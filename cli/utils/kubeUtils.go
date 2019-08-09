package utils

import (
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

// IsKubectlAvailable checks whether kubectl is available
func IsKubectlAvailable() (bool, error) {

	_, err := keptnutils.ExecuteCommand("kubectl", []string{})
	if err != nil {
		return false, err
	}
	return true, nil
}

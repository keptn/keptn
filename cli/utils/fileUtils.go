package utils

import (
	"os"
	"os/user"
)

const keptnFolderName = ".keptn"

// GetKeptnDirectory returns a path, which is used to store logs and possibly creds
func GetKeptnDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	keptnDir := usr.HomeDir + string(os.PathSeparator) + keptnFolderName + string(os.PathSeparator)

	if _, err := os.Stat(keptnDir); os.IsNotExist(err) {
		err := os.MkdirAll(keptnDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return keptnDir, nil
}

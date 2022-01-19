package common

import (
	"fmt"
	"os"
)

func GetProjectConfigPath(project string) string {
	return fmt.Sprintf("%s/%s", GetConfigDir(), project)
}

func GetProjectMetadataFilePath(project string) string {
	return fmt.Sprintf("%s/%s", GetProjectConfigPath(project), "metadata.yaml")
}

func GetServiceConfigPath(project, service string) string {
	return fmt.Sprintf("%s/%s", GetProjectConfigPath(project), service)
}

func ensureDirectoryExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		err := os.MkdirAll(path, 0700)
		return err
	}
	return nil
}

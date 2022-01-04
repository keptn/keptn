package common

import "fmt"

func GetProjectConfigPath(project string) string {
	return fmt.Sprintf("%s/%s", GetConfigDir(), project)
}

func GetProjectMetadataFilePath(project string) string {
	return fmt.Sprintf("%s/%s", GetProjectConfigPath(project), "metadata.yaml")
}

func GetServiceConfigPath(project, service string) string {
	return fmt.Sprintf("%s/%s", GetProjectConfigPath(project), service)
}

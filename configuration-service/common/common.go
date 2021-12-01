package common

import (
	"fmt"
	"github.com/keptn/keptn/configuration-service/config"
)

const ConfigurationServiceName = "configuration-service"

const ProjectDoesNotExistErrorMsg = "Project does not exist"
const StageDoesNotExistErrorMsg = "Stage does not exist"
const ServiceDoesNotExistErrorMsg = "Service does not exist"

const CannotCheckOutBranchErrorMsg = "Could not check out branch"
const GitURLNotFound = "Repository not found"
const HostNotFound = "host"
const GitError = "exit status 128"
const WrongToken = "access token"

const CannotAddResourceErrorMsg = "Could not add resource"
const CannotUpdateResourceErrorMsg = "Could not update resource"

const StageDirectoryName = "keptn-stages"
const ServiceDirectoryName = "keptn-services"

func GetProjectConfigPath(project string) string {
	return fmt.Sprintf("%s/%s", config.ConfigDir, project)
}

func GetStageConfigPath(project, stage string) string {
	return fmt.Sprintf("%s/%s/%s", GetProjectConfigPath(project), StageDirectoryName, stage)
}

func GetServiceConfigPath(project, stage, service string) string {
	return fmt.Sprintf("%s/%s/%s", GetStageConfigPath(project, stage), ServiceDirectoryName, service)
}

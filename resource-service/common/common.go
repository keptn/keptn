package common

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/config"
)

const ErrReadOnly = "read-only filesystem"

const ConfigurationServiceName = "configuration-service"

const ProjectDoesNotExistErrorMsg = "Project does not exist"
const StageDoesNotExistErrorMsg = "Stage does not exist"
const ServiceDoesNotExistErrorMsg = "Service does not exist"
const InvalidContextErrorMsg = "git context is invalid"
const InvalidCredentialsErrorMsg = "git credentials are invalid"

func GetProjectConfigPath(project string) string {
	return fmt.Sprintf("%s/%s", config.ConfigDir, project)
}

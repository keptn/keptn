package common

import (
	"fmt"
	"github.com/keptn/keptn/configuration-service/config"
)

const ErrReadOnly = "read-only filesystem"

const ConfigurationServiceName = "configuration-service"

const ProjectDoesNotExistErrorMsg = "Project does not exist"
const StageDoesNotExistErrorMsg = "Stage does not exist"
const ServiceDoesNotExistErrorMsg = "Service does not exist"

func GetProjectConfigPath(project string) string {
	return fmt.Sprintf("%s/%s", config.ConfigDir, project)
}

package common

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/config"
)

const gitKeptnUserDefault = "keptn"
const gitKeptnEmailDefault = "keptn@keptn.sh"
const gitKeptnUserEnvVar = "GIT_KEPTN_USER"
const gitKeptnEmailEnvVar = "GIT_KEPTN_EMAIL"

func GetProjectConfigPath(project string) string {
	return fmt.Sprintf("%s/%s", config.ConfigDir, project)
}

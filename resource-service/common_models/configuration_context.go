package common_models

import "github.com/keptn/keptn/resource-service/models"

type ConfigurationContextParams struct {
	Project                 models.Project
	Stage                   *models.Stage
	Service                 *models.Service
	GitContext              GitContext
	CheckConfigDirAvailable bool
}

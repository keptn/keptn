package lib

import (
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"
	"strings"
)

type ControlPlane struct {
	UniformHandler *api.UniformHandler
	CurrentID      string
	EnvConfig      config.EnvConfig
}

func (c *ControlPlane) Register() error {
	logger.Info("Registering integration")
	data := c.getRegistrationDataFromEnv()
	id, err := c.UniformHandler.RegisterIntegration(data)
	if err != nil {
		return err
	}
	c.CurrentID = id
	return nil
}

func (c *ControlPlane) Unregister() error {
	logger.Info("Unregistering integration")
	err := c.UniformHandler.UnregisterIntegration(c.CurrentID)
	if err != nil {
		return err
	}
	c.CurrentID = ""
	return nil
}

func (c *ControlPlane) getRegistrationDataFromEnv() models.Integration {
	return models.Integration{
		Name: c.EnvConfig.K8sPodName,
		MetaData: models.MetaData{
			Hostname:           c.EnvConfig.K8sNodeName,
			DeploymentName:     c.EnvConfig.K8sDeploymentName,
			IntegrationVersion: c.EnvConfig.Version,
			DistributorVersion: c.EnvConfig.DistributorVersion,
			Status:             "",
			Location:           c.EnvConfig.Location,
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      c.EnvConfig.K8sNamespace,
				PodName:        c.EnvConfig.K8sPodName,
				DeploymentName: c.EnvConfig.K8sDeploymentName,
			},
		},
		Subscription: models.Subscription{
			Topics: strings.Split(c.EnvConfig.PubSubTopic, ","),
			Status: "UP",
			Filter: models.SubscriptionFilter{
				Project: c.EnvConfig.ProjectFilter,
				Stage:   c.EnvConfig.StageFilter,
				Service: c.EnvConfig.ServiceFilter,
			},
		},
	}

}

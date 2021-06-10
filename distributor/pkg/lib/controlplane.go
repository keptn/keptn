package lib

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type ControlPlane struct {
	UniformHandler *api.UniformHandler
	EnvConfig      config.EnvConfig
	currentID      string
	mux            sync.Mutex
}

func (c *ControlPlane) Register() (string, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	logger.Info("Registering integration")
	data := c.getRegistrationDataFromEnv()
	id, err := c.UniformHandler.RegisterIntegration(data)
	if err != nil {
		return "", err
	}
	c.currentID = id
	return c.currentID, nil
}

func (c *ControlPlane) Unregister() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	logger.Info("Unregistering integration")
	if c.currentID == "" {
		return fmt.Errorf("tried to unrigster without being registered first")
	}
	err := c.UniformHandler.UnregisterIntegration(c.currentID)
	if err != nil {
		return err
	}
	c.currentID = ""
	return nil
}

func (c *ControlPlane) getRegistrationDataFromEnv() models.Integration {
	var topics []string
	if c.EnvConfig.PubSubTopic == "" {
		topics = []string{}
	} else {
		topics = strings.Split(c.EnvConfig.PubSubTopic, ",")
	}
	return models.Integration{
		Name: c.EnvConfig.K8sPodName,
		MetaData: models.MetaData{
			Hostname:           c.EnvConfig.K8sNodeName,
			DeploymentName:     c.EnvConfig.K8sDeploymentName,
			IntegrationVersion: c.EnvConfig.Version,
			DistributorVersion: c.EnvConfig.DistributorVersion,
			Location:           c.EnvConfig.Location,
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      c.EnvConfig.K8sNamespace,
				PodName:        c.EnvConfig.K8sPodName,
				DeploymentName: c.EnvConfig.K8sDeploymentName,
			},
		},
		Subscription: models.Subscription{
			Topics: topics,
			Filter: models.SubscriptionFilter{
				Project: c.EnvConfig.ProjectFilter,
				Stage:   c.EnvConfig.StageFilter,
				Service: c.EnvConfig.ServiceFilter,
			},
		},
	}

}

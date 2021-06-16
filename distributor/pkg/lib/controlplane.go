package lib

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/distributor/pkg/config"
	"strings"
	"sync"
)

type ControlPlane struct {
	UniformHandler *api.UniformHandler
	ConnectionType config.ConnectionType
	EnvConfig      config.EnvConfig
	currentID      string
	mux            sync.Mutex
}

func NewControlPlane(connectionType config.ConnectionType, env config.EnvConfig) *ControlPlane {
	var uniformHandler *api.UniformHandler
	if connectionType == config.ConnectionTypeHTTP {
		uniformHandler = api.NewAuthenticatedUniformHandler(env.KeptnAPIEndpoint+"/controlPlane", env.KeptnAPIToken, "x-token", nil, "http")
	} else {
		uniformHandler = api.NewUniformHandler(config.DefaultShipyardControllerBaseURL)
	}

	return &ControlPlane{
		UniformHandler: uniformHandler,
		ConnectionType: connectionType,
		EnvConfig:      env,
		currentID:      "",
		mux:            sync.Mutex{},
	}
}

func (c *ControlPlane) Register() (string, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
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
	if c.currentID == "" {
		return fmt.Errorf("tried to unregister integration without being registered first")
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

	var location string
	if c.EnvConfig.Location == "" {
		location = config.ConnectionTypeToLocation[c.ConnectionType]
	} else {
		location = c.EnvConfig.Location
	}
	return models.Integration{
		Name: c.EnvConfig.K8sDeploymentName,
		MetaData: models.MetaData{
			Hostname:           c.EnvConfig.K8sNodeName,
			IntegrationVersion: c.EnvConfig.Version,
			DistributorVersion: c.EnvConfig.DistributorVersion,
			Location:           location,
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

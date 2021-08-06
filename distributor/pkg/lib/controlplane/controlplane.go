package controlplane

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
	sync.Mutex
	uniformHandler  *api.UniformHandler
	currentID       string
	integrationData models.Integration
}

func NewControlPlane(uniformHandler *api.UniformHandler, integrationData models.Integration) *ControlPlane {
	return &ControlPlane{
		uniformHandler:  uniformHandler,
		integrationData: integrationData,
	}
}

func (c *ControlPlane) Register() (string, error) {
	c.Lock()
	defer c.Unlock()
	id, err := c.uniformHandler.RegisterIntegration(c.integrationData)
	if err != nil {
		return "", err
	}
	c.currentID = id
	return c.currentID, nil
}

func (c *ControlPlane) Unregister() error {
	c.Lock()
	defer c.Unlock()
	if c.currentID == "" {
		return fmt.Errorf("tried to unregister integration without being registered first")
	}
	err := c.uniformHandler.UnregisterIntegration(c.currentID)
	if err != nil {
		return err
	}
	c.currentID = ""
	return nil
}

func CreateRegistrationData(connectionType config.ConnectionType, env config.EnvConfig) models.Integration {
	var topics []string
	if env.PubSubTopic == "" {
		topics = []string{}
	} else {
		topics = strings.Split(env.PubSubTopic, ",")
	}

	var location string
	if env.Location == "" {
		location = config.ConnectionTypeToLocation[connectionType]
	} else {
		location = env.Location
	}

	var stageFilter []string
	if env.StageFilter == "" {
		stageFilter = []string{}
	} else {
		stageFilter = strings.Split(env.StageFilter, ",")
	}

	var serviceFilter []string
	if env.ServiceFilter == "" {
		serviceFilter = []string{}
	} else {
		serviceFilter = strings.Split(env.ServiceFilter, ",")
	}

	var projectFilter []string
	if env.ProjectFilter == "" {
		projectFilter = []string{}
	} else {
		projectFilter = strings.Split(env.ProjectFilter, ",")
	}

	if env.K8sNodeName == "" {
		logger.Warn("K8S_NODE_NAME is not set. Using default value: 'keptn-node'")
		env.K8sNodeName = "keptn-node"
	}

	return models.Integration{
		Name: env.K8sDeploymentName,
		MetaData: models.MetaData{
			Hostname:           env.K8sNodeName,
			IntegrationVersion: env.Version,
			DistributorVersion: env.DistributorVersion,
			Location:           location,
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      env.K8sNamespace,
				PodName:        env.K8sPodName,
				DeploymentName: env.K8sDeploymentName,
			},
		},
		Subscriptions: []models.Subscription{{
			Topics: topics,
			Filter: models.SubscriptionFilter{
				Project:  env.ProjectFilter,
				Projects: projectFilter,
				Stage:    env.StageFilter,
				Stages:   stageFilter,
				Service:  env.ServiceFilter,
				Services: serviceFilter,
			},
		},
		},
	}
}

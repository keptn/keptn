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
	uniformHandler *api.UniformHandler
	connectionType config.ConnectionType
	env            config.EnvConfig
	currentID      string
}

func NewControlPlane(uniformHandler *api.UniformHandler, connectionType config.ConnectionType, env config.EnvConfig) *ControlPlane {
	return &ControlPlane{
		uniformHandler: uniformHandler,
		connectionType: connectionType,
		env:            env,
	}
}

func (c *ControlPlane) Ping() (*models.Integration, error) {
	return c.uniformHandler.Ping(c.currentID)
}

func (c *ControlPlane) Register() (string, error) {
	c.Lock()
	defer c.Unlock()
	id, err := c.uniformHandler.RegisterIntegration(c.createRegistrationData())
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

func (c *ControlPlane) createRegistrationData() models.Integration {

	var location string
	if c.env.Location == "" {
		location = config.ConnectionTypeToLocation[c.connectionType]
	} else {
		location = c.env.Location
	}

	var stageFilter []string
	if c.env.StageFilter == "" {
		stageFilter = []string{}
	} else {
		stageFilter = strings.Split(c.env.StageFilter, ",")
	}

	var serviceFilter []string
	if c.env.ServiceFilter == "" {
		serviceFilter = []string{}
	} else {
		serviceFilter = strings.Split(c.env.ServiceFilter, ",")
	}

	var projectFilter []string
	if c.env.ProjectFilter == "" {
		projectFilter = []string{}
	} else {
		projectFilter = strings.Split(c.env.ProjectFilter, ",")
	}

	if c.env.K8sNodeName == "" {
		logger.Warn("K8S_NODE_NAME is not set. Using default value: 'keptn-node'")
		c.env.K8sNodeName = "keptn-node"
	}

	//create subscription
	topics := []string{}
	if c.env.PubSubTopic == "" {
		topics = []string{}
	} else {
		topics = strings.Split(c.env.PubSubTopic, ",")
	}
	var subscriptions []models.EventSubscription
	for _, t := range topics {
		ts := models.EventSubscription{
			Event: t,
			Filter: models.EventSubscriptionFilter{
				Projects: projectFilter,
				Stages:   stageFilter,
				Services: serviceFilter,
			},
		}
		subscriptions = append(subscriptions, ts)
	}

	return models.Integration{
		Name: c.env.K8sDeploymentName,
		MetaData: models.MetaData{
			Hostname:           c.env.K8sNodeName,
			IntegrationVersion: c.env.Version,
			DistributorVersion: c.env.DistributorVersion,
			Location:           location,
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      c.env.K8sNamespace,
				PodName:        c.env.K8sPodName,
				DeploymentName: c.env.K8sDeploymentName,
			},
		},
		Subscriptions: subscriptions,
	}
}

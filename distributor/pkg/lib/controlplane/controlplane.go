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

type IControlPlane interface {
	Ping() (*models.Integration, error)
	Register() (string, error)
	Unregister() error
}

type ControlPlane struct {
	sync.Mutex
	uniformHandler *api.UniformHandler
	connectionType config.ConnectionType
	currentID      string
}

func NewControlPlane(uniformHandler *api.UniformHandler, connectionType config.ConnectionType) *ControlPlane {
	return &ControlPlane{
		uniformHandler: uniformHandler,
		connectionType: connectionType,
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
	if config.Global.Location == "" {
		location = config.ConnectionTypeToLocation[c.connectionType]
	} else {
		location = config.Global.Location
	}

	var stageFilter []string
	if config.Global.StageFilter == "" {
		stageFilter = []string{}
	} else {
		stageFilter = strings.Split(config.Global.StageFilter, ",")
	}

	var serviceFilter []string
	if config.Global.ServiceFilter == "" {
		serviceFilter = []string{}
	} else {
		serviceFilter = strings.Split(config.Global.ServiceFilter, ",")
	}

	var projectFilter []string
	if config.Global.ProjectFilter == "" {
		projectFilter = []string{}
	} else {
		projectFilter = strings.Split(config.Global.ProjectFilter, ",")
	}

	if config.Global.K8sNodeName == "" {
		logger.Warn("K8S_NODE_NAME is not set. Using default value: 'keptn-node'")
		config.Global.K8sNodeName = "keptn-node"
	}

	//create subscription
	topics := []string{}
	if config.Global.PubSubTopic == "" {
		topics = []string{}
	} else {
		topics = strings.Split(config.Global.PubSubTopic, ",")
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
		Name: config.Global.K8sDeploymentName,
		MetaData: models.MetaData{
			Hostname:           config.Global.K8sNodeName,
			IntegrationVersion: config.Global.Version,
			DistributorVersion: config.Global.DistributorVersion,
			Location:           location,
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      config.Global.K8sNamespace,
				PodName:        config.Global.K8sPodName,
				DeploymentName: config.Global.K8sDeploymentName,
			},
		},
		Subscriptions: subscriptions,
	}
}

package lib

import (
	"fmt"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/distributor/pkg/config"
	logger "github.com/sirupsen/logrus"
)

type ControlPlane struct {
	UniformHandler *api.UniformHandler
	CurrentID      string
	EnvConfig      config.EnvConfig
}

func (c *ControlPlane) Register() error {
	logger.Info("Registering integration")
	fmt.Println(c.EnvConfig.KeptnAPIEndpoint)
	//integration := models.Integration{
	//	Name: "",
	//	MetaData: models.MetaData{
	//		Hostname:           "",
	//		DeploymentName:     "",
	//		IntegrationVersion: "",
	//		DistributorVersion: "",
	//		Status:             "",
	//		Location:           "",
	//		KubernetesMetaData: models.KubernetesMetaData{},
	//	},
	//	Subscription: models.Subscription{
	//		Topics: nil,
	//		Status: "",
	//		Filter: models.SubscriptionFilter{
	//			Project: "",
	//			Stage:   "",
	//			Service: "",
	//		},
	//	},
	//}
	//registerIntegration, err := c.UniformHandler.RegisterIntegration(integration)
	//if err != nil {
	//	return fmt.Errorf(*err.Message)
	//}
	//c.currentID = registerIntegration
	return nil
}

func (c *ControlPlane) Unregister() error {
	logger.Info("Unregistering integration")
	return nil
}

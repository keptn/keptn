package handlers

import (
	keptnmongoutils "github.com/keptn/go-utils/pkg/common/mongoutils"
	log "github.com/sirupsen/logrus"
)

const (
	eventsCollectionName = "keptnUnmappedEvents"
	serviceName          = "mongodb-datastore"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	K8SDeploymentVersion string `envconfig:"K8S_DEPLOYMENT_VERSION" default:""`
	K8SPodName           string `envconfig:"K8S_POD_NAME" default:""`
	K8SNamespace         string `envconfig:"K8S_NAMESPACE" default:""`
	K8SNodeName          string `envconfig:"K8S_NODE_NAME" default:""`
	LogLevel             string `envconfig:"LOG_LEVEL" default:"info"`
}

func GetMongoDBConnectionString() (string, string, error) {
	return keptnmongoutils.GetMongoConnectionStringFromEnv()
}

func (env envConfig) ConfigLog() {
	log.SetLevel(log.InfoLevel)
	if env.LogLevel != "" {
		logLevel, err := log.ParseLevel(env.LogLevel)
		if err != nil {
			log.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			log.SetLevel(logLevel)
		}
	}

}

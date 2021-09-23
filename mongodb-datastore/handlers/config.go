package handlers

import (
	keptnmongoutils "github.com/keptn/go-utils/pkg/common/mongoutils"
)

const (
	eventsCollectionName = "keptnUnmappedEvents"
	serviceName          = "mongodb-datastore"
)

func GetMongoDBConnectionString() (string, string, error) {
	return keptnmongoutils.GetMongoConnectionStringFromEnv()
}

package handlers

import (
	"fmt"
	"os"
)

var (
	mongoDBHost       = os.Getenv("MONGODB_HOST")
	mongoDBName       = os.Getenv("MONGODB_DATABASE")
	mongoDBUser       = os.Getenv("MONGODB_USER")
	mongoDBPassword   = os.Getenv("MONGODB_PASSWORD")
	mongoDBReplicaSet = os.Getenv("MONGODB_REPLICASET")
	mongoDBConnection = fmt.Sprintf("mongodb://%s:%s@%s/%s?replicaSet=%s", mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBName, mongoDBReplicaSet)
)

const (
	eventsCollectionName = "keptnUnmappedEvents"
	serviceName          = "mongodb-datastore"
)

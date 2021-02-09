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
	mongoDBConnection = fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBName)
)

const (
	eventsCollectionName = "keptnUnmappedEvents"
	serviceName          = "mongodb-datastore"
)

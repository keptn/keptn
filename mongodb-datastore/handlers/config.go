package handlers

import (
	"fmt"
	"os"
)

var mongoDBHost = os.Getenv("MONGODB_HOST")
var mongoDBName = os.Getenv("MONGODB_DATABASE")
var mongoDBUser = os.Getenv("MONGODB_USER")
var mongoDBPassword = os.Getenv("MONGODB_PASSWORD")

var mongoDBConnection = fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBName)

const eventsCollectionName = "keptnUnmappedEvents"
const serviceName = "mongodb-datastore"

package handlers

import "os"

var mongoDBConnection = os.Getenv("MONGO_DB_CONNECTION_STRING")
var mongoDBName = os.Getenv("MONGO_DB_NAME")

const eventsCollectionName = "events"

const serviceName = "mongodb-datastore"

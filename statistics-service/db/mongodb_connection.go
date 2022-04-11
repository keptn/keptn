package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	keptnmongoutils "github.com/keptn/go-utils/pkg/common/mongoutils"

	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var databaseName string

var mutex = &sync.Mutex{}

const clientCreationFailed = "failed to create mongo client: %v"
const clientConnectionFailed = "failed to create mongo client: %v"

// MongoDBConnection takes care of establishing a connection to the mongodb
type MongoDBConnection struct {
	Client *mongo.Client
}

// EnsureDBConnection makes sure a connection to the mongodb is established
func (m *MongoDBConnection) EnsureDBConnection() error {
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	// attention: not calling the cancel() function likely causes memory leaks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if m.Client == nil {
		logger.Info("No MongoDB client has been initialized yet. Creating a new one.")
		return m.connectMongoDBClient()
	} else if err = m.Client.Ping(ctx, nil); err != nil {
		logger.Info("MongoDB client lost connection. Attempt reconnect.")
		ctxDisconnect, cancelDisconnect := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancelDisconnect()
		err2 := m.Client.Disconnect(ctxDisconnect)
		if err2 != nil {
			logger.Errorf("failed to disconnect client from MongoDB: %v", err2)
		}
		return m.connectMongoDBClient()
	}
	return nil
}

func (m *MongoDBConnection) connectMongoDBClient() error {
	var err error
	connectionString, dbName, err := keptnmongoutils.GetMongoConnectionStringFromEnv()
	if err != nil {
		logger.Errorf(clientCreationFailed, err)
		return fmt.Errorf(clientCreationFailed, err)
	}
	databaseName = dbName
	m.Client, err = mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		logger.Errorf(clientCreationFailed, err)
		return fmt.Errorf(clientCreationFailed, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = m.Client.Connect(ctx)
	if err != nil {
		logger.Infof(clientConnectionFailed, err)
		return fmt.Errorf(clientConnectionFailed, err)
	}
	return nil
}

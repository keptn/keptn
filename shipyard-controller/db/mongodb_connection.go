package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	keptnmongoutils "github.com/keptn/go-utils/pkg/common/mongoutils"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mutex = &sync.Mutex{}

var mongoDBConnectionInstance *MongoDBConnection

var mongoConnectionOnce sync.Once

// MongoDBConnection takes care of establishing a connection to the mongodb
type MongoDBConnection struct {
	Client *mongo.Client
}

func GetMongoDBConnectionInstance() *MongoDBConnection {
	mongoConnectionOnce.Do(func() {
		mongoDBConnectionInstance = &MongoDBConnection{}
	})
	return mongoDBConnectionInstance
}

func getDatabaseName() string {
	return os.Getenv("MONGODB_DATABASE")
}

// EnsureDBConnection makes sure a connection to the mongodb is established
func (m *MongoDBConnection) EnsureDBConnection() error {
	var err error
	// attention: not calling the cancel() function likely causes memory leaks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if m.Client == nil {
		logger.Info("No MongoDB client has been initialized yet. Creating a new one.")
		return m.connectMongoDBClient()
	} else if err = m.Client.Ping(ctx, nil); err != nil {
		logger.Info("MongoDB client lost connection. Attempt reconnect.")
		err2 := m.Client.Disconnect(ctx)
		if err2 != nil {
			logger.Errorf("failed to disconnect client from MongoDB: %v", err2)
		}
		return m.connectMongoDBClient()
	}
	return nil
}

func (m *MongoDBConnection) connectMongoDBClient() error {
	connectionString, _, err := keptnmongoutils.GetMongoConnectionStringFromEnv()
	if err != nil {
		logger.Errorf("failed to create mongo client: %v", err)
		return fmt.Errorf("failed to create mongo client: %v", err)
	}
	m.Client, err = mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		logger.Errorf("failed to create mongo client: %v", err)
		return fmt.Errorf("failed to create mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = m.Client.Connect(ctx)
	if err != nil {
		logger.Infof("failed to connect client to MongoDB: %v", err)
		return fmt.Errorf("failed to connect client to MongoDB: %v", err)
	}
	return nil
}

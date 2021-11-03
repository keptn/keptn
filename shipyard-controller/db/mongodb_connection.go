package db

import (
	"context"
	"fmt"
	keptnmongoutils "github.com/keptn/go-utils/pkg/common/mongoutils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
	"time"
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
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	// attention: not calling the cancel() function likely causes memory leaks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if m.Client == nil {
		fmt.Println("No MongoDB client has been initialized yet. Creating a new one.")
		return m.connectMongoDBClient()
	} else if err = m.Client.Ping(ctx, nil); err != nil {
		fmt.Println("MongoDB client lost connection. Attempt reconnect.")
		return m.connectMongoDBClient()
	}
	return nil
}

func (m *MongoDBConnection) connectMongoDBClient() error {
	var err error

	connectionString, _, err := keptnmongoutils.GetMongoConnectionStringFromEnv()
	if err != nil {
		return fmt.Errorf("failed to create mongo client: %v", err)
	}
	m.Client, err = mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		err := fmt.Errorf("failed to create mongo client: %v", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = m.Client.Connect(ctx)
	if err != nil {
		err := fmt.Errorf("failed to connect client to MongoDB: %v", err)
		return err
	}
	return nil
}

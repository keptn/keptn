package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDBHost = os.Getenv("MONGODB_HOST")
var databaseName = os.Getenv("MONGO_DB_NAME")
var mongoDBUser = os.Getenv("MONGODB_USER")
var mongoDBPassword = os.Getenv("MONGODB_PASSWORD")
var mutex = &sync.Mutex{}

var mongoDBConnection = fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, databaseName)

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
	m.Client, err = mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
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

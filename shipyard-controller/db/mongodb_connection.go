package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBConnection struct {
	Client *mongo.Client
}

func (m *MongoDBConnection) EnsureDBConnection() error {
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	if m.Client == nil {
		fmt.Println("No MongoDB client has been initialized yet. Creating a new one.")
		return m.connectMongoDBClient()
	} else if err = m.Client.Ping(context.TODO(), nil); err != nil {
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

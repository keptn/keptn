package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/keptn/go-utils/pkg/common/strutils"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mutex = &sync.Mutex{}

var mongoDBConnectionInstance *MongoDBConnection

var mongoConnectionOnce sync.Once

const clientCreationFailed = "failed to create mongo client: %v"
const clientConnectionFailed = "failed to connect client to MongoDB: %v"

// MongoDBConnection takes care of establishing a connection to the mongodb
type MongoDBConnection struct {
	client *mongo.Client
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

func (m *MongoDBConnection) GetClient() (*mongo.Client, error) {
	err := m.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	return m.client, nil
}

// EnsureDBConnection makes sure a connection to the mongodb is established
func (m *MongoDBConnection) EnsureDBConnection() error {
	// Mutex is neccessary for not creating multiple clients after restarting the pod
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	// attention: not calling the cancel() function likely causes memory leaks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if m.client == nil {
		logger.Debug("No MongoDB client has not been initialized yet. Creating a new one.")
		return m.connectMongoDBClient()
	} else if err = m.client.Ping(ctx, nil); err != nil {
		logger.Warn("MongoDB client lost connection. Attempting reconnect.")
		ctxDisconnect, cancelDisconnect := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancelDisconnect()
		err2 := m.client.Disconnect(ctxDisconnect)
		if err2 != nil {
			logger.Errorf("failed to disconnect client from MongoDB: %v", err2)
		}
		m.client = nil
		return m.connectMongoDBClient()
	}
	return nil
}

func (m *MongoDBConnection) connectMongoDBClient() error {
	connectionString, _, err := GetMongoConnectionStringFromFile()
	if err != nil {
		logger.Errorf(clientCreationFailed, err)
		return fmt.Errorf(clientCreationFailed, err)
	}
	clientOptions := options.Client()
	clientOptions = clientOptions.ApplyURI(connectionString)
	clientOptions = clientOptions.SetConnectTimeout(30 * time.Second)
	m.client, err = mongo.NewClient(clientOptions)
	if err != nil {
		logger.Errorf(clientCreationFailed, err)
		return fmt.Errorf(clientCreationFailed, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = m.client.Connect(ctx)
	if err != nil {
		logger.Errorf(clientConnectionFailed, err)
		return fmt.Errorf(clientConnectionFailed, err)
	}
	if err = m.client.Ping(ctx, nil); err != nil {
		logger.Errorf(clientConnectionFailed, err)
		return fmt.Errorf(clientConnectionFailed, err)
	}

	logger.Info("Successfully connected to MongoDB")
	return nil
}

const mongoUser = "mongodb-user/"
const mongoPwd = "mongodb-passwords/"
const mongoExtCon = "external_connection_string/"

// mongodb://<MONGODB_USER>:<MONGODB_PASSWORD>@MONGODB_HOST>/<MONGODB_DATABASE>
func GetMongoConnectionStringFromFile() (string, string, error) {
	mongoDBName := os.Getenv("MONGODB_DATABASE")
	if mongoDBName == "" {
		return "", "", errors.New("env var 'MONGODB_DATABASE' env var must be set")
	}
	configdir := os.Getenv("MONGO_CONFIG_DIR")

	if externalConnectionString := readSecret(configdir + mongoExtCon); externalConnectionString != "" {
		return externalConnectionString, mongoDBName, nil
	}
	mongoDBHost := os.Getenv("MONGODB_HOST")
	mongoDBUser := readSecret(configdir + mongoUser)
	mongoDBPassword := readSecret(configdir + mongoPwd)

	if !strutils.AllSet(mongoDBHost, mongoDBUser, mongoDBPassword) {
		return "", "", errors.New("could not construct mongodb connection string: env vars 'MONGODB_HOST', 'MONGODB_USER' and 'MONGODB_PASSWORD' have to be set")
	}
	return fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBName), mongoDBName, nil
}

func readSecret(file string) string {
	body, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}
	return string(body)
}

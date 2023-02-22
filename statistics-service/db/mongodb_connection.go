package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/keptn/go-utils/pkg/common/strutils"

	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var databaseName string

var mutex = &sync.Mutex{}

const clientCreationFailed = "failed to create mongo client: %v"
const clientConnectionFailed = "failed to connect client to MongoDB: %v"

// MongoDBConnection takes care of establishing a connection to the mongodb
type MongoDBConnection struct {
	Client *mongo.Client
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
	if m.Client == nil {
		logger.Info("No MongoDB client has been initialized yet. Creating a new one.")
		return m.connectMongoDBClient()
	} else if err = m.Client.Ping(ctx, nil); err != nil {
		logger.Info("MongoDB client lost connection. Attempting reconnect.")
		ctxDisconnect, cancelDisconnect := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancelDisconnect()
		err2 := m.Client.Disconnect(ctxDisconnect)
		if err2 != nil {
			logger.Errorf("failed to disconnect client from MongoDB: %v", err2)
		}
		m.Client = nil
		return m.connectMongoDBClient()
	}
	return nil
}

func (m *MongoDBConnection) connectMongoDBClient() error {
	var err error
	connectionString, dbName, err := GetMongoConnectionStringFromFile()
	if err != nil {
		logger.Errorf(clientCreationFailed, err)
		return fmt.Errorf(clientCreationFailed, err)
	}
	databaseName = dbName
	clientOptions := options.Client()
	clientOptions = clientOptions.ApplyURI(connectionString)
	clientOptions = clientOptions.SetConnectTimeout(30 * time.Second)
	m.Client, err = mongo.NewClient(clientOptions)
	if err != nil {
		logger.Errorf(clientCreationFailed, err)
		return fmt.Errorf(clientCreationFailed, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = m.Client.Connect(ctx)
	if err != nil {
		logger.Errorf(clientConnectionFailed, err)
		return fmt.Errorf(clientConnectionFailed, err)
	}
	if err = m.Client.Ping(ctx, nil); err != nil {
		logger.Errorf(clientConnectionFailed, err)
		return fmt.Errorf(clientConnectionFailed, err)
	}

	logger.Info("Successfully connected to MongoDB")
	return nil
}

const mongoUser = "/mongodb-user"
const mongoPwd = "/mongodb-passwords"
const mongoExtCon = "/external_connection_string"

// mongodb://<MONGODB_USER>:<MONGODB_PASSWORD>@MONGODB_HOST>/<MONGODB_DATABASE>
func GetMongoConnectionStringFromFile() (string, string, error) {
	mongoDBName := os.Getenv("MONGODB_DATABASE")
	if mongoDBName == "" {
		return "", "", errors.New("env var 'MONGODB_DATABASE' env var must be set")
	}
	configdir := os.Getenv("MONGO_CONFIG_DIR")

	if externalConnectionString := getFromEnvOrFile("MONGODB_EXTERNAL_CONNECTION_STRING", configdir+mongoExtCon); externalConnectionString != "" {
		return externalConnectionString, mongoDBName, nil
	}
	mongoDBHost := os.Getenv("MONGODB_HOST")
	mongoDBUser := getFromEnvOrFile("MONGODB_USER", configdir+mongoUser)
	mongoDBPassword := getFromEnvOrFile("MONGODB_PASSWORD", configdir+mongoPwd)

	if !strutils.AllSet(mongoDBHost, mongoDBUser, mongoDBPassword) {
		return "", "", errors.New("could not construct mongodb connection string: env vars 'MONGODB_HOST', 'MONGODB_USER' and 'MONGODB_PASSWORD' have to be set")
	}
	return fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBName), mongoDBName, nil
}

func readSecret(file string) string {
	body, err := os.ReadFile(file)
	if err != nil {
		logger.Fatalf("unable to read secret: %v", err)
	}
	return string(body)
}

func getFromEnvOrFile(env string, path string) string {
	if _, err := os.Stat(path); err == nil {
		return readSecret(path)
	}
	return os.Getenv(env)
}

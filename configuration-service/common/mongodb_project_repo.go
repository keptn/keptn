package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/keptn/keptn/configuration-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDBHost = os.Getenv("MONGODB_HOST")
var databaseName = os.Getenv("MONGO_DB_NAME")
var mongoDBUser = os.Getenv("MONGODB_USER")
var mongoDBPassword = os.Getenv("MONGODB_PASSWORD")

var mongoDBConnection = fmt.Sprintf("mongodb://%s:%s@%s/%s", mongoDBUser, mongoDBPassword, mongoDBHost, databaseName)

const projectsCollectionName = "keptnProjectsMV"

type MongoDBProjectRepo struct {
	Client *mongo.Client
}

func (mdbrepo *MongoDBProjectRepo) CreateProject(project *models.ExpandedProject) error {
	err := mdbrepo.ensureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(project)

	projectCollection := mdbrepo.getProjectsCollection()
	_, err = projectCollection.InsertOne(ctx, prjInterface)
	if err != nil {
		fmt.Println("Could not create project " + project.ProjectName + ": " + err.Error())
	}
	return nil
}

func (mdbrepo *MongoDBProjectRepo) GetProject(projectName string) (*models.ExpandedProject, error) {
	err := mdbrepo.ensureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mdbrepo.getProjectsCollection()
	result := projectCollection.FindOne(ctx, bson.M{"projectName": projectName})
	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	projectResult := &models.ExpandedProject{}
	err = result.Decode(projectResult)
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not cast %v to *models.Project\n", result))
		return nil, err
	}
	return projectResult, nil
}
func (mdbrepo *MongoDBProjectRepo) GetProjects() ([]*models.ExpandedProject, error) {
	result := []*models.ExpandedProject{}
	err := mdbrepo.ensureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mdbrepo.getProjectsCollection()
	cursor, err := projectCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error retrieving projects from mongoDB: " + err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		projectResult := &models.ExpandedProject{}
		err := cursor.Decode(projectResult)
		if err != nil {
			fmt.Println("Could not cast to *models.Project")
		}
		result = append(result, projectResult)
	}

	return result, nil
}
func (mdbrepo *MongoDBProjectRepo) UpdateProject(project *models.ExpandedProject) error {
	err := mdbrepo.ensureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(project)
	projectCollection := mdbrepo.getProjectsCollection()
	_, err = projectCollection.ReplaceOne(ctx, bson.M{"projectName": project.ProjectName}, prjInterface)
	if err != nil {
		fmt.Println("Could not update project " + project.ProjectName + ": " + err.Error())
		return err
	}
	return nil
}
func (mdbrepo *MongoDBProjectRepo) DeleteProject(projectName string) error {
	err := mdbrepo.ensureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mdbrepo.getProjectsCollection()
	_, err = projectCollection.DeleteMany(ctx, bson.M{"projectName": projectName})
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not delete project %s : %s\n", projectName, err.Error()))
		return err
	}
	fmt.Println("Deleted project " + projectName)
	return nil
}

func (mdbrepo *MongoDBProjectRepo) ensureDBConnection() error {
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	if mdbrepo.Client == nil {
		fmt.Println("No MongoDB client has been initialized yet. Creating a new one.")
		return mdbrepo.connectMongoDBClient()
	} else if err = mdbrepo.Client.Ping(context.TODO(), nil); err != nil {
		fmt.Println("MongoDB client lost connection. Attempt reconnect.")
		return mdbrepo.connectMongoDBClient()
	}
	return nil
}

func (mdbrepo *MongoDBProjectRepo) connectMongoDBClient() error {
	var err error
	mdbrepo.Client, err = mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		err := fmt.Errorf("failed to create mongo client: %v", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = mdbrepo.Client.Connect(ctx)
	if err != nil {
		err := fmt.Errorf("failed to connect client to MongoDB: %v", err)
		return err
	}
	return nil
}

func (mdbrepo *MongoDBProjectRepo) getProjectsCollection() *mongo.Collection {
	projectCollection := mdbrepo.Client.Database(databaseName).Collection(projectsCollectionName)
	return projectCollection
}

func transformProjectToInterface(prj *models.ExpandedProject) interface{} {
	// marshall and unmarshall again because for some reason the json tags of the golang struct of the project type are not considered
	marshal, _ := json.Marshal(prj)
	var prjInterface interface{}
	json.Unmarshal(marshal, &prjInterface)
	return prjInterface
}

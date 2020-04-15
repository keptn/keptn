package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/keptn/keptn/configuration-service/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var mongoDBConnection = "mongodb://user:password@localhost:27017/keptn" //os.Getenv("MONGO_DB_CONNECTION_STRING")
var databaseName = "keptn"                                              //os.Getenv("MONGO_DB_NAME")
const projectsCollectionName = "keptnProjectsMV"

var instance *mongoDBMaterializedView

type mongoDBMaterializedView struct {
	Client *mongo.Client
}

func GetMongoDBMaterializedView() *mongoDBMaterializedView {
	if instance == nil {
		instance = &mongoDBMaterializedView{}
	}
	return instance
}

func (mmv *mongoDBMaterializedView) CreateProject(prj *models.Project) error {
	existingProject, err := mmv.GetProject(prj.ProjectName)
	if existingProject != nil {
		return nil
	}
	err = mmv.createProject(prj)
	if err != nil {
		return err
	}
	return nil
}

func (mmv *mongoDBMaterializedView) UpdateShipyard(projectName string, shipyardContent string) error {
	existingProject, err := mmv.GetProject(projectName)
	if err != nil {
		return err
	}

	existingProject.Shipyard = shipyardContent

	return mmv.updateProject(projectName, existingProject)
}

func (mmv *mongoDBMaterializedView) GetProjects() (*models.ExpandedProjects, error) {
	result := &models.ExpandedProjects{Projects: []*models.ExpandedProject{}}
	err := mmv.ensureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mmv.getProjectsCollection()
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
		result.Projects = append(result.Projects, projectResult)
	}

	return result, nil
}

func (mmv *mongoDBMaterializedView) getProjectsCollection() *mongo.Collection {
	projectCollection := mmv.Client.Database(databaseName).Collection(projectsCollectionName)
	return projectCollection
}

func (mmv *mongoDBMaterializedView) GetProject(project string) (*models.ExpandedProject, error) {
	err := mmv.ensureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mmv.getProjectsCollection()
	result := projectCollection.FindOne(ctx, bson.M{"projectName": project})
	projectResult := &models.ExpandedProject{}
	err = result.Decode(projectResult)
	if err != nil {
		fmt.Sprintf("Could not cast %v to *models.Project\n", result)
		return nil, err
	}
	return projectResult, nil
}

func (mmv *mongoDBMaterializedView) DeleteProject(project string) error {
	err := mmv.ensureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := mmv.getProjectsCollection()
	_, err = projectCollection.DeleteMany(ctx, bson.M{"projectName": project})
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not delete project %s : %s\n", project, err.Error()))
		return err
	}
	fmt.Println("Deleted project " + project)
	return nil
}

func (mmv *mongoDBMaterializedView) CreateStage(project string, stage string) error {
	fmt.Println("Adding stage " + stage + " to project " + project)
	prj, err := mmv.GetProject(project)

	if err != nil {
		fmt.Sprintf("Could not add stage %s to project %s : %s\n", stage, project, err.Error())
		return err
	}

	stageAlreadyExists := false
	for _, stg := range prj.Stages {
		if stg.StageName == stage {
			stageAlreadyExists = true
			break
		}
	}

	if stageAlreadyExists {
		fmt.Println("Stage " + stage + " already exists in project " + project)
		return nil
	}

	prj.Stages = append(prj.Stages, &models.ExpandedStage{
		Services:  []*models.ExpandedService{},
		StageName: stage,
	})

	err = mmv.updateProject(project, prj)
	if err != nil {
		return err
	}

	fmt.Println("Added stage " + stage + " to project " + project)
	return nil
}

func (mmv *mongoDBMaterializedView) createProject(prj *models.Project) error {
	expandedProject := &models.ExpandedProject{
		CreationDate: time.Now().String(),
		GitRemoteURI: prj.GitRemoteURI,
		GitUser:      prj.GitUser,
		ProjectName:  prj.ProjectName,
		Shipyard:     "",
		Stages:       nil,
	}

	err := mmv.ensureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(expandedProject)

	projectCollection := mmv.getProjectsCollection()
	_, err = projectCollection.InsertOne(ctx, prjInterface)
	if err != nil {
		fmt.Println("Could not create project " + prj.ProjectName + ": " + err.Error())
	}
	return nil
}

func transformProjectToInterface(prj *models.ExpandedProject) interface{} {
	// marshall and unmarshall again because for some reason the json tags of the golang struct of the project type are not considered
	marshal, _ := json.Marshal(prj)
	var prjInterface interface{}
	json.Unmarshal(marshal, &prjInterface)
	return prjInterface
}

func (mmv *mongoDBMaterializedView) updateProject(project string, prj *models.ExpandedProject) error {
	err := mmv.ensureDBConnection()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prjInterface := transformProjectToInterface(prj)
	projectCollection := mmv.getProjectsCollection()
	_, err = projectCollection.ReplaceOne(ctx, bson.M{"projectName": project}, prjInterface)
	if err != nil {
		fmt.Println("Could not update project " + project + ": " + err.Error())
	}
	return nil
}

func (mmv *mongoDBMaterializedView) DeleteStage(project string, stage string) error {
	fmt.Println("Deleting stage " + stage + " from project " + project)
	prj, err := mmv.GetProject(project)

	if err != nil {
		fmt.Sprintf("Could not delete stage %s from project %s : %s\n", stage, project, err.Error())
		return err
	}

	stageIndex := 0

	for idx, stg := range prj.Stages {
		if stg.StageName == stage {
			stageIndex = idx
			break
		}
	}

	copy(prj.Stages[stageIndex:], prj.Stages[stageIndex+1:])
	prj.Stages[len(prj.Stages)-1] = nil
	prj.Stages = prj.Stages[:len(prj.Stages)-1]

	err = mmv.updateProject(project, prj)
	return nil
}

func (mmv *mongoDBMaterializedView) CreateService(project string, stage string, service string) error {
	existingProject, err := mmv.GetProject(project)
	if err != nil {
		fmt.Println("Could not add service " + service + " to stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
		return err
	}

	for _, stg := range existingProject.Stages {
		if stg.StageName == stage {
			for _, svc := range stg.Services {
				if svc.ServiceName == service {
					break
				}
			}
			stg.Services = append(stg.Services, &models.ExpandedService{
				CreationDate:  time.Now().String(),
				DeployedImage: "",
				ServiceName:   service,
			})
			fmt.Println("Adding " + service + " to stage " + stage + " in project " + project + " in database")
			err := mmv.updateProject(project, existingProject)
			if err != nil {
				fmt.Println("Could not add service " + service + " to stage " + stage + " in project " + project + ". Could not update project: " + err.Error())
			}
			break
		}
	}
	fmt.Println("Service " + service + " already exists in stage " + stage + " in project " + project)
	return nil
}

func (mmv *mongoDBMaterializedView) DeleteService(project string, stage string, service string) error {
	return nil
}

func (mmv *mongoDBMaterializedView) ensureDBConnection() error {
	mutex.Lock()
	defer mutex.Unlock()
	var err error
	if mmv.Client == nil {
		fmt.Println("No MongoDB client has been initialized yet. Creating a new one.")
		return mmv.connectMongoDBClient()
	} else if err = mmv.Client.Ping(context.TODO(), nil); err != nil {
		fmt.Println("MongoDB client lost connection. Attempt reconnect.")
		return mmv.connectMongoDBClient()
	}
	return nil
}

func (mmv *mongoDBMaterializedView) connectMongoDBClient() error {
	var err error
	mmv.Client, err = mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		err := fmt.Errorf("failed to create mongo client: %v", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = mmv.Client.Connect(ctx)
	if err != nil {
		err := fmt.Errorf("failed to connect client to MongoDB: %v", err)
		return err
	}
	return nil
}

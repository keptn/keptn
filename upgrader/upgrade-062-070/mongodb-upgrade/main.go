package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"

	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
)

const defaultMongoDBConnectionString = "mongodb://user:password@mongodb.keptn-datastore.svc.cluster.local:27017/keptn"
const defaultConfigurationServiceURL = "configuration-service.keptn.svc.cluster.local:8080"

var client *mongo.Client

var projectCollections map[string]*mongo.Collection

func main() {
	var mongoDBConnectionString string
	var configurationServiceURL string
	if len(os.Args) > 1 {
		mongoDBConnectionString = os.Args[1]
	} else {
		mongoDBConnectionString = defaultMongoDBConnectionString
	}

	if len(os.Args) > 2 {
		configurationServiceURL = os.Args[2]
	} else {
		configurationServiceURL = defaultConfigurationServiceURL
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnectionString))
	if err != nil {
		fmt.Printf("failed to create mongo client: %v\n", err)
		os.Exit(1)
	}
	ctx := context.TODO()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Printf("failed to connect client to MongoDB: %v\n", err)
		os.Exit(1)
	}

	projectsMV := []*keptnapimodels.Project{}

	projectHandler := keptnapi.NewProjectHandler(configurationServiceURL)
	stageHandler := keptnapi.NewStageHandler(configurationServiceURL)
	serviceHandler := keptnapi.NewServiceHandler(configurationServiceURL)

	allProjects, err := projectHandler.GetAllProjects()
	if err != nil {
		fmt.Println("failed to retrieve projects from configuration service: " + err.Error())
		os.Exit(1)
	}

	for _, prj := range allProjects {
		stages, err := stageHandler.GetAllStages(prj.ProjectName)
		if err != nil {
			fmt.Println("failed to retrieve stages of project " + prj.ProjectName + " from configuration service: " + err.Error())
			os.Exit(1)
		}

		for _, stg := range stages {
			services, err := serviceHandler.GetAllServices(prj.ProjectName, stg.StageName)
			if err != nil {
				fmt.Println("failed to retrieve services of project " + prj.ProjectName + " from configuration service: " + err.Error())
				os.Exit(1)
			}

			stg.Services = services
		}

		prj.Stages = stages
		projectsMV = append(projectsMV, prj)
	}

	eventsCollection := client.Database("keptn").Collection("events")
	contextToProjectCollection := client.Database("keptn").Collection("contextToProject")
	projectCollections = map[string]*mongo.Collection{}

	// get all events from events collection
	sortOptions := options.Find().SetSort(bson.D{{"time", 1}})

	cursor, err := eventsCollection.Find(ctx, bson.D{}, sortOptions)
	if err != nil {
		fmt.Printf("failed to retrieve events from mongodb: %v\n", err)
		os.Exit(1)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc bson.M
		err := cursor.Decode(&doc)
		if err != nil {
			fmt.Printf("failed to decode event %v\n", err)
			os.Exit(1)
		}

		keptnContext, ok := doc["shkeptncontext"].(string)
		if !ok || keptnContext == "" {
			fmt.Printf("Cannot migrate event because object does not have the expected structure.\n")
			continue
		}

		data, ok := doc["data"].(bson.M)
		if !ok || data == nil {
			fmt.Printf("Cannot migrate event because object does not have the expected structure.\n")
			continue
		}
		project, ok := data["project"].(string)
		if !ok || project == "" {
			fmt.Printf("Cannot migrate event because object does not have the expected structure.\n")
			continue
		}

		if projectCollections[project] == nil {
			projectCollections[project] = client.Database("keptn").Collection(project)
		}
		_, err = projectCollections[project].InsertOne(ctx, doc)
		if err != nil {
			writeErr, ok := err.(mongo.WriteException)
			if ok {
				if len(writeErr.WriteErrors) > 0 && writeErr.WriteErrors[0].Code == 11000 { // 11000 = duplicate key error
					fmt.Printf("Event %v already exists in collection\n", doc)
				}
			} else {
				fmt.Printf("Could not store event %v\n", doc)
			}
		}
		fmt.Printf("Inserted event %v into collection %s\n", doc, project)
		_, err = contextToProjectCollection.InsertOne(ctx, bson.M{"_id": keptnContext, "shkeptncontext": keptnContext, "project": project})
		if err != nil {
			writeErr, ok := err.(mongo.WriteException)
			if ok {
				if len(writeErr.WriteErrors) > 0 && writeErr.WriteErrors[0].Code == 11000 { // 11000 = duplicate key error
					fmt.Printf("Mapping %s -> %s  already exists in collection\n", keptnContext, project)
				}
			} else {
				fmt.Printf("Could not store mapping %s -> %s\n", keptnContext, project)
			}
		}
		fmt.Printf("Inserted mapping %s -> %s\n", keptnContext, project)
	}

}

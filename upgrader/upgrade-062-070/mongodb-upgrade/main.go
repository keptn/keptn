package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
)

const defaultMongoDBConnectionString = "mongodb://user:password@mongodb.keptn-datastore.svc.cluster.local:27017/keptn"
const defaultConfigurationServiceURL = "configuration-service.keptn.svc.cluster.local:8080"

const materializedViewCollection = "keptnProjectsMV"

var client *mongo.Client

var projectCollections map[string]*mongo.Collection

var mongoDBConnectionString string
var configurationServiceURL string
var skipWrite = true

func main() {

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

	allProjects, err := getAllProjects()
	if err != nil {
		fmt.Println("failed to retrieve projects from configuration service: " + err.Error())
		os.Exit(1)
	}

	for _, prj := range allProjects {
		stages, err := getAllStages(prj.ProjectName)
		if err != nil {
			fmt.Println("failed to retrieve stages of project " + prj.ProjectName + " from configuration service: " + err.Error())
			os.Exit(1)
		}

		for _, stg := range stages {
			services, err := getAllServices(prj.ProjectName, stg.StageName)
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

		updateLastEventOfService(projectsMV, doc)

		fmt.Printf("Inserting event %v into collection %s\n", doc, project)
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

	fmt.Println(fmt.Printf("Projects Materialized View:\n%v", projectsMV))

	mvCollection := client.Database("keptn").Collection(materializedViewCollection)
	for _, projectMV := range projectsMV {

		existingProject := mvCollection.FindOne(ctx, bson.M{"projectName": projectMV.ProjectName})

		if existingProject.Err() == nil {
			// project already exists - must not recreate it
			fmt.Println("Project " + projectMV.ProjectName + " already exists in MV table.")
			continue
		}

		projectInterface, _ := transformProjectToInterface(projectMV)
		_, err = mvCollection.InsertOne(ctx, projectInterface)
		if err != nil {
			writeErr, ok := err.(mongo.WriteException)
			if ok {
				if len(writeErr.WriteErrors) > 0 && writeErr.WriteErrors[0].Code == 11000 { // 11000 = duplicate key error
					fmt.Println("Project already exists in collection")
				}
			} else {
				fmt.Printf("Could not store materialized view of project %s\n", projectMV.ProjectName)
			}
		}
	}
}

func transformProjectToInterface(prj *keptnapimodels.Project) (interface{}, error) {
	data, err := json.Marshal(prj)
	if err != nil {
		err := fmt.Errorf("failed to marshal event: %v", err)
		return nil, err
	}

	var projectInterface interface{}
	err = json.Unmarshal(data, &projectInterface)
	if err != nil {
		err := fmt.Errorf("failed to unmarshal event: %v", err)
		return nil, err
	}
	return projectInterface, nil
}

func updateLastEventOfService(projectsMV []*keptnapimodels.Project, event bson.M) {
	data, ok := event["data"].(bson.M)
	if !ok || data == nil {
		fmt.Printf("Cannot assign event to service because object does not have the expected structure.\n")
		return
	}
	eventType, ok := event["type"].(string)
	if !ok || eventType == "" {
		fmt.Printf("Cannot assign event to service because object does not have the expected structure.\n")
		return
	}
	project, ok := data["project"].(string)
	if !ok || project == "" {
		fmt.Printf("Cannot assign event to service because object does not have the expected structure.\n")
		return
	}
	stage, ok := data["stage"].(string)
	if !ok || stage == "" {
		fmt.Printf("Cannot assign event to service because object does not have the expected structure.\n")
		return
	}
	service, ok := data["service"].(string)
	if !ok || service == "" {
		fmt.Printf("Cannot assign event to service because object does not have the expected structure.\n")
		return
	}
	keptnContext, ok := event["shkeptncontext"].(string)
	if !ok || keptnContext == "" {
		fmt.Printf("Cannot assign event to service because object does not have the expected structure.\n")
		return
	}
	eventID, ok := event["id"].(string)
	if !ok || eventID == "" {
		fmt.Printf("Cannot assign event to service because object does not have the expected structure.\n")
		return
	}
	time, ok := event["time"].(string)
	if !ok || time == "" {
		fmt.Printf("Cannot assign event to service because object does not have the expected structure.\n")
		return
	}

	for _, prj := range projectsMV {
		if prj.ProjectName == project {
			for _, stg := range prj.Stages {
				if stg.StageName == stage {
					for _, svc := range stg.Services {
						if svc.ServiceName == service {
							if svc.LastEventTypes == nil {
								svc.LastEventTypes = map[string]keptnapimodels.EventContextInfo{}
							}
							svc.LastEventTypes[eventType] = keptnapimodels.EventContextInfo{
								EventID:      eventID,
								KeptnContext: keptnContext,
								Time:         time,
							}

							if eventType == keptn.DeploymentFinishedEventType {
								image, ok := data["image"].(string)
								if !ok || image == "" {
									fmt.Printf("Cannot update deployed image of service because object does not have the expected structure.\n")
									continue
								}
								tag, ok := data["tag"].(string)
								if !ok || tag == "" {
									fmt.Printf("Cannot update deployed image of service because object does not have the expected structure.\n")
									continue
								}
								svc.DeployedImage = image + ":" + tag
							}
						}
					}
				}
			}
		}
	}
}

func getAllProjects() ([]*keptnapimodels.Project, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	projects := []*keptnapimodels.Project{}

	nextPageKey := ""

	for {
		url, err := url.Parse("http://" + configurationServiceURL + "/v1/project/")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")

		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 200 {
			var received keptnapimodels.Projects
			err = json.Unmarshal(body, &received)
			if err != nil {
				return nil, err
			}
			projects = append(projects, received.Projects...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey
		} else {
			var respErr keptnapimodels.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, err
			}
			return nil, errors.New("Response Error Code: " + string(respErr.Code) + " Message: " + *respErr.Message)
		}
	}

	return projects, nil
}

func getAllStages(project string) ([]*keptnapimodels.Stage, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	stages := []*keptnapimodels.Stage{}

	nextPageKey := ""
	for {
		url, err := url.Parse("http://" + configurationServiceURL + "/v1/project/" + project + "/stage")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			var received keptnapimodels.Stages
			err = json.Unmarshal(body, &received)
			if err != nil {
				return nil, err
			}
			stages = append(stages, received.Stages...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey
		} else {
			var respErr keptnapimodels.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, err
			}
			return nil, errors.New("Response Error Code: " + string(respErr.Code) + " Message: " + *respErr.Message)
		}
	}
	return stages, nil
}

func getAllServices(project string, stage string) ([]*keptnapimodels.Service, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	services := []*keptnapimodels.Service{}

	nextPageKey := ""

	for {
		url, err := url.Parse("http://" + configurationServiceURL + "/v1/project/" + project + "/stage/" + stage + "/service")
		if err != nil {
			return nil, err
		}
		q := url.Query()
		if nextPageKey != "" {
			q.Set("nextPageKey", nextPageKey)
			url.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", url.String(), nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 200 {
			var received keptnapimodels.Services
			err = json.Unmarshal(body, &received)
			if err != nil {
				return nil, err
			}
			services = append(services, received.Services...)

			if received.NextPageKey == "" || received.NextPageKey == "0" {
				break
			}
			nextPageKey = received.NextPageKey
		} else {
			var respErr keptnapimodels.Error
			err = json.Unmarshal(body, &respErr)
			if err != nil {
				return nil, err
			}
			return nil, errors.New("Response Error Code: " + string(respErr.Code) + " Message: " + *respErr.Message)
		}
	}

	return services, nil
}

package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	keptn "github.com/keptn/go-utils/pkg/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
)

const defaultMongoDBTargetConnectionString = "mongodb://user:password@mongodb:27017/keptn"
const defaultMongoDBSourceConnectionString = "mongodb://user:password@mongodb.keptn-datastore:27017/keptn"
const defaultConfigurationServiceURL = "configuration-service.keptn:8080"

const rootEventCollectionSuffix = "-rootEvents"

const materializedViewCollection = "keptnProjectsMV"

const projectsMVFile = "projects-mv.json"

type ProjectsMV struct {
	Projects []*keptnapimodels.Project `json:"projects"`
}

var projectCollections map[string]*mongo.Collection

var mongoDBSourceConnectionString string
var mongoDBTargetConnectionString string
var configurationServiceURL string
var action string
var skipWrite = true

func main() {

	if len(os.Args) > 1 {
		mongoDBSourceConnectionString = os.Args[1]
	} else {
		mongoDBSourceConnectionString = defaultMongoDBSourceConnectionString
	}

	if len(os.Args) > 2 {
		mongoDBTargetConnectionString = os.Args[2]
	} else {
		mongoDBTargetConnectionString = defaultMongoDBTargetConnectionString
	}

	if len(os.Args) > 3 {
		configurationServiceURL = os.Args[3]
	} else {
		configurationServiceURL = defaultConfigurationServiceURL
	}

	const storeProjectsMVAction = "store-projects-mv"
	if len(os.Args) > 4 {
		action = os.Args[4]
	} else {
		action = ""
	}

	if action == storeProjectsMVAction {
		projectsMVJSONObject := &ProjectsMV{}
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
		projectsMVJSONObject.Projects = projectsMV
		projectsMVJSONString, err := json.MarshalIndent(projectsMVJSONObject, "", " ")
		if err != nil {
			fmt.Println("Could not store projects MV: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("Projects MV before migration:")
		fmt.Println(string(projectsMVJSONString))

		_ = ioutil.WriteFile(projectsMVFile, projectsMVJSONString, 0644)
		return
	}

	jsonFile, err := os.Open(projectsMVFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	projectsMVJSONObject := &ProjectsMV{}
	err = json.Unmarshal(byteValue, projectsMVJSONObject)
	if err != nil {
		fmt.Println("Could not read projects MV: " + err.Error())
	}

	projectsMV := projectsMVJSONObject.Projects

	sourceClient, sourceCtx, err := createMongoDBClient(mongoDBSourceConnectionString)
	if err != nil {
		fmt.Printf("failed to create mongo sourceClient: %v\n", err)
		os.Exit(1)
	}

	targetClient, targetCtx, err := createMongoDBClient(mongoDBTargetConnectionString)
	if err != nil {
		fmt.Printf("failed to create mongo sourceClient: %v\n", err)
		os.Exit(1)
	}

	eventsSourceCollection := sourceClient.Database("keptn").Collection("events")
	contextToProjectTargetCollection := targetClient.Database("keptn").Collection("contextToProject")
	projectCollections = map[string]*mongo.Collection{}

	// get all events from events collection
	sortOptions := options.Find().SetSort(bson.D{{"time", 1}})

	cursor, err := eventsSourceCollection.Find(sourceCtx, bson.D{}, sortOptions)
	if err != nil {
		fmt.Printf("failed to retrieve events from mongodb: %v\n", err)
		os.Exit(1)
	}
	defer cursor.Close(sourceCtx)
	for cursor.Next(sourceCtx) {
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
			projectCollections[project] = targetClient.Database("keptn").Collection(project)
		}
		_, err = projectCollections[project].InsertOne(targetCtx, doc)
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
		_, err = contextToProjectTargetCollection.InsertOne(targetCtx, bson.M{"_id": keptnContext, "shkeptncontext": keptnContext, "project": project})
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

		err = storeRootEvent(targetClient, project, targetCtx, doc)
		if err != nil {
			fmt.Println("Could not store root event: " + err.Error())
		}
	}

	fmt.Println(fmt.Printf("Projects Materialized View:\n%v", projectsMV))

	mvCollection := targetClient.Database("keptn").Collection(materializedViewCollection)
	for _, projectMV := range projectsMV {

		existingProject := mvCollection.FindOne(targetCtx, bson.M{"projectName": projectMV.ProjectName})

		if existingProject.Err() == nil {
			// project already exists - must not recreate it
			fmt.Println("Project " + projectMV.ProjectName + " already exists in MV table.")
			continue
		}

		projectInterface, _ := transformProjectToInterface(projectMV)
		_, err = mvCollection.InsertOne(targetCtx, projectInterface)
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

func createMongoDBClient(url string) (*mongo.Client, context.Context, error) {
	sourceClient, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, nil, err
	}
	ctx := context.TODO()

	err = sourceClient.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}
	return sourceClient, ctx, err
}

func transformEventToInterface(event interface{}) (interface{}, error) {
	data, err := json.Marshal(event)
	if err != nil {
		err := fmt.Errorf("failed to marshal event: %v", err)
		return nil, err
	}

	var eventInterface interface{}
	err = json.Unmarshal(data, &eventInterface)
	if err != nil {
		err := fmt.Errorf("failed to unmarshal event: %v", err)
		return nil, err
	}
	return eventInterface, nil
}

func storeRootEvent(client *mongo.Client, collectionName string, ctx context.Context, event bson.M) error {

	keptnContext := event["shkeptncontext"].(string)
	event["_id"] = keptnContext

	rootEventsForProjectCollection := client.Database("keptn").Collection(collectionName + rootEventCollectionSuffix)

	result := rootEventsForProjectCollection.FindOne(ctx, bson.M{"shkeptncontext": keptnContext})

	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		eventInterface, err := transformEventToInterface(event)
		if err != nil {
			fmt.Println("Could not transform root event to interface: " + err.Error())
			return err
		}

		_, err = rootEventsForProjectCollection.InsertOne(ctx, eventInterface)
		if err != nil {
			err := fmt.Errorf("Failed to store root event for KeptnContext "+keptnContext+": %v", err.Error())
			fmt.Println(err.Error())
			return err
		}
		fmt.Println("Stored root event for KeptnContext: " + keptnContext)
	} else if result.Err() != nil {
		// found an already stored root event => check if incoming event precedes the already existing event
		// if yes, then the new event will be the new root event for this context
		existingEvent := &keptnapimodels.KeptnContextExtendedCE{}

		err := result.Decode(existingEvent)
		if err != nil {
			fmt.Println("Could not decode existing root event: " + err.Error())
			return err
		}

		if time.Time(existingEvent.Time).After(event["time"].(time.Time)) {
			fmt.Println("Replacing root event for KeptnContext: " + keptnContext)
			_, err := rootEventsForProjectCollection.DeleteOne(ctx, bson.M{"_id": existingEvent.ID})
			if err != nil {
				fmt.Println("Could not delete previous root event: " + err.Error())
				return err
			}
			eventInterface, err := transformEventToInterface(event)
			if err != nil {
				fmt.Println("Could not transform root event to interface: " + err.Error())
				return err
			}

			_, err = rootEventsForProjectCollection.InsertOne(ctx, eventInterface)
			if err != nil {
				err := fmt.Errorf("Failed to store root event for KeptnContext "+keptnContext+": %v", err.Error())
				fmt.Println(err.Error())
				return err
			}
			fmt.Println("Stored root event for KeptnContext: " + keptnContext)
		}
	}
	fmt.Println("Root event for KeptnContext " + keptnContext + " already exists in collection")
	return nil
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

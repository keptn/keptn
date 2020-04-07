package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const mongoDBURL = "mongodb://user:password@localhost:27017/keptn"

var client *mongo.Client

var projectCollections map[string]*mongo.Collection

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBURL))
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

	eventsCollection := client.Database("keptn").Collection("events")
	contextToProjectCollection := client.Database("keptn").Collection("contextToProject")
	projectCollections = map[string]*mongo.Collection{}
	// get all events from events collection
	cursor, err := eventsCollection.Find(ctx, bson.D{})
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

		if doc["data"] == nil || project == "" || doc["shkeptncontext"].(string) == "" {
			fmt.Printf("Cannot migrate event because no project has been detected.")
			continue
		}
		if projectCollections[project] == nil {
			projectCollections[project] = client.Database("keptn").Collection(project)
		}
		_, err = projectCollections[project].InsertOne(ctx, doc)
		if err != nil {
			fmt.Printf("Could not store event %v\n", doc)
		}
		fmt.Printf("Inserted event %v into collection %s\n", doc, project)
		_, err = contextToProjectCollection.InsertOne(ctx, bson.M{"_id": keptnContext, "shkeptncontext": keptnContext, "project": project})
		if err != nil {
			fmt.Printf("Could not store mapping %s -> %s\n", keptnContext, project)
		}
		fmt.Printf("Inserted mapping %s -> %s\n", keptnContext, project)
	}

}

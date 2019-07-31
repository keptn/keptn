package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/keptn/mongodb-datastore/restapi/operations/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveLog to datastore
func SaveLog(body []*logs.SaveLogParamsBodyItems0) (err error) {
	fmt.Println("save log to datastore")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		log.Fatalln("error creating client: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	collection := client.Database(mongoDBName).Collection(logsCollectionName)

	for _, l := range body {

		res, err := collection.InsertOne(ctx, l)
		if err != nil {
			log.Fatalln("error inserting: ", err.Error())
		}
		fmt.Println("insertedID: ", res.InsertedID)
	}

	return err

}

// GetLogs returns logs
func GetLogs() (res []*logs.GetLogsOKBodyItems0, err error) {
	fmt.Println("get logs from datastore")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBConnection))
	if err != nil {
		log.Fatalln("error creating client: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	collection := client.Database(mongoDBName).Collection(logsCollectionName)

	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatalln("error finding elements in collections: ", err.Error())
	}

	var resultLogs []*logs.GetLogsOKBodyItems0
	for cur.Next(ctx) {
		var result logs.GetLogsOKBodyItems0
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		resultLogs = append(resultLogs, &result)
		//fmt.Println(result)
	}
	return resultLogs, nil
}

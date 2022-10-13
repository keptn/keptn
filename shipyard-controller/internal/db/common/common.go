package common

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

func CloseCursor(ctx context.Context, cur *mongo.Cursor) {
	if cur == nil {
		return
	}
	if err := cur.Close(ctx); err != nil {
		log.Errorf("Could not close cursor: %v", err)
	}
}

func EncodeKey(key string) string {
	encodedKey := strings.ReplaceAll(strings.ReplaceAll(key, "~", "~t"), ".", "~p")
	return encodedKey
}
func DecodeKey(key string) string {
	decodedKey := strings.ReplaceAll(strings.ReplaceAll(key, "~p", "."), "~t", "~")
	return decodedKey
}

func ToInterface(item interface{}) (interface{}, error) {
	// marshall and unmarshall again because for some reason the json tags of the golang struct of the project type are not considered
	marshal, _ := json.Marshal(item)
	var prjInterface interface{}
	err := json.Unmarshal(marshal, &prjInterface)
	if err != nil {
		return nil, err
	}
	return prjInterface, nil
}

package db

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func closeCursor(ctx context.Context, cur *mongo.Cursor) {
	if cur == nil {
		return
	}
	if err := cur.Close(ctx); err != nil {
		log.Errorf("Could not close cursor: %v", err)
	}
}

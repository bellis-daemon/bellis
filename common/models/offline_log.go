package models

import (
	"time"

	"github.com/bellis-daemon/bellis/common/models/index"
	"github.com/bellis-daemon/bellis/common/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OfflineLog struct {
	ID             primitive.ObjectID `bson:"_id"`
	EntityID       primitive.ObjectID `bson:"EntityID"`
	EnvoyType      string             `bson:"EnvoyType"`
	OfflineTime    time.Time          `bson:"OfflineTime"`
	OfflineMessage string             `bson:"OfflineMessage"`
	OnlineTime     time.Time          `bson:"OnlineTime"`
	SentryLogs     []SentryLog        `bson:"SentryLogs"`
}

type SentryLog struct {
	SentryName   string    `bson:"SentryName"`
	SentryTime   time.Time `bson:"SentryTime"`
	ErrorMessage string    `bson:"ErrorMessage"`
}

func init() {
	index.RegistrerIndex(&storage.COfflineLog, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "EntityID", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "OfflineTime", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "OnlineTime", Value: 1},
			},
		},
	})
}

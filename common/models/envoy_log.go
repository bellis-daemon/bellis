package models

import (
	"time"

	"github.com/bellis-daemon/bellis/common/models/index"
	"github.com/bellis-daemon/bellis/common/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnvoyLog struct {
	ID             primitive.ObjectID `json:"ID" bson:"_id"`
	SendTime       time.Time          `json:"SendTime" bson:"SendTime"`
	Success        bool               `json:"Success" bson:"Success"`
	FailedMessage  string             `json:"FaildMessage" bson:"FailedMessage"`
	OfflineLogID   primitive.ObjectID `json:"OfflineLogID" bson:"OfflineLogID"`
	PolicyType     string             `json:"PolicyType" bson:"PolicyType"`
	PolicySnapShot bson.M             `json:"PolicySnapShot" bson:"PolicySnapShot"`
}

func init() {
	index.RegistrerIndex(&storage.CEnvoyLog, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "OfflineLogID", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "SendTime", Value: 1},
			},
		},
	})
}

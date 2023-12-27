package models

import (
	"time"

	"github.com/bellis-daemon/bellis/common/models/index"
	"github.com/bellis-daemon/bellis/common/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserLoginLog struct {
	ID         primitive.ObjectID `json:"ID" bson:"_id"`
	UserID     primitive.ObjectID `json:"UserID" bson:"UserID"`
	LoginTime  time.Time          `json:"LoginTime" bson:"LoginTime"`
	Location   string             `json:"Location" bson:"Location"`
	Device     string             `json:"Device" bson:"Device"`
	DeviceType string             `json:"DeviceType" bson:"DeviceType"`
}

func init() {
	index.RegistrerIndex(&storage.CUserLoginLog, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "UserID", Value: 1},
			},
		},
	})
}

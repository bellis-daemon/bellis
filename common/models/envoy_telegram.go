package models

import (
	"github.com/bellis-daemon/bellis/common/models/index"
	"github.com/bellis-daemon/bellis/common/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnvoyTelegram struct {
	EnvoyHeader
	ID     primitive.ObjectID `json:"ID" bson:"_id"`
	ChatID int64              `json:"ChatID" bson:"ChatID"`
}

func init() {
	index.RegistrerIndex(&storage.CEnvoyTelegram, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "ChatID", Value: 1},
			},
		},
	})
}

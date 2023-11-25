package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnvoyTelegram struct {
	ID     primitive.ObjectID `json:"ID" bson:"_id"`
	ChatId int64              `json:"ChatId" bson:"ChatId"`
}

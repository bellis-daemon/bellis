package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnvoyEmail struct {
	EnvoyHeader
	ID      primitive.ObjectID `json:"ID" bson:"_id"`
	Address string             `json:"Address" bson:"Address"`
}

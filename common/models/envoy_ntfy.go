package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnvoyNtfy struct {
	ID primitive.ObjectID `json:"ID" bson:"_id"`
}

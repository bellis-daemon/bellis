package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnvoyGotify struct {
	EnvoyHeader `json:"EnvoyHeader" bson:"EnvoyHeader"`
	ID    primitive.ObjectID `json:"ID" bson:"_id"`
	URL   string             `json:"URL" bson:"URL"`
	Token string             `json:"Token" bson:"Token"`
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnvoyWebhook struct {
	ID       primitive.ObjectID `json:"ID" bson:"_id"`
	URL      string             `json:"URL" bson:"URL"`
	Insecure bool               `json:"Insecure" bson:"Insecure"`
}

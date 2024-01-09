package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnvoyWebhook struct {
	EnvoyHeader `json:"EnvoyHeader" bson:"EnvoyHeader"`
	ID            primitive.ObjectID `json:"ID" bson:"_id"`
	URL           string             `json:"URL" bson:"URL"`
	Insecure      bool               `json:"Insecure" bson:"Insecure"`     // true: HTTP, false: HTTPS
	AuthMethod    string             `json:"AuthMethod" bson:"AuthMethod"` // "None" "Basic" "Bearer"
	BasicUsername string             `json:"BasicUsername" bson:"BasicUsername"`
	BasicPassword string             `json:"BasicPassword" bson:"BasicPassword"`
	BearerToken   string             `json:"BearerToken" bson:"BearerToken"`
}

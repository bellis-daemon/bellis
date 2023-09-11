package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Application struct {
	ID          primitive.ObjectID       `json:"ID" bson:"_id"`
	Name        string                   `json:"Name" bson:"Name"`
	Description string                   `json:"Description" bson:"Description"`
	UserID      primitive.ObjectID       `json:"UserID" bson:"UserID"`
	CreatedAt   time.Time                `json:"CreatedAt" bson:"CreatedAt"`
	Scheme      string                   `json:"Scheme" bson:"Scheme"`
	Active      bool                     `json:"Active" bson:"Active"`
	Public      ApplicationPublicOptions `json:"Public" bson:"Public"`
	Options     bson.M                   `json:"options" bson:"options"`
}

type ApplicationPublicOptions struct {
	Threshold   int      `json:"Threshold" bson:"Threshold"`
	TriggerList []string `json:"TriggerList" bson:"TriggerList"`
}

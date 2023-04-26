package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Application struct {
	ID          primitive.ObjectID `json:"ID" bson:"_id"`
	Name        string             `json:"Name" bson:"Name"`
	Description string             `json:"Description" bson:"Description"`
	UserID      primitive.ObjectID `json:"UserID" bson:"UserID"`
	CreatedAt   time.Time          `json:"CreatedAt" bson:"CreatedAt"`
	SchemeID    int                `json:"SchemeID" bson:"SchemeID"`
	Active      bool               `json:"Active" bson:"Active"`
	Options     bson.M             `json:"Options" bson:"Options"`
}

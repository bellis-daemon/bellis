package models

import (
	"context"
	"time"

	"github.com/bellis-daemon/bellis/common/models/index"
	"github.com/bellis-daemon/bellis/common/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Multiplier  uint     `json:"Multiplier" bson:"Multiplier"`
	Threshold   int      `json:"Threshold" bson:"Threshold"`
	TriggerList []string `json:"TriggerList" bson:"TriggerList"`
}

func (this *Application) User() (*User, error) {
	var user User
	err := storage.CUser.FindOne(context.Background(), bson.M{
		"_id": this.UserID,
	}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func init() {
	index.RegistrerIndex(&storage.CEntity, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "UserID", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "Name", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "CreatedAt", Value: 1},
			},
		},
	})
}

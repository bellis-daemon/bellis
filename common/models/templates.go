package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Timezone string

const DefaultTimezone = "Asia/Shanghai"

func (this Timezone) Location() *time.Location {
	s := string(this)
	if this == "" {
		s = DefaultTimezone
	}
	loc, _ := time.LoadLocation(s)
	return loc
}

type EnvoyHeader struct {
	UserID       primitive.ObjectID `json:"UserID" bson:"UserID"`
	CreatedAt    time.Time          `json:"CreatedAt" bson:"CreatedAt"`
	OfflineAlert bool               `json:"OfflineAlert" bson:"OfflineAlert"`
	Sensitive    int                `json:"Sensitive" bson:"Sensitive"`
}

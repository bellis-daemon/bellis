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
	UserID    primitive.ObjectID `json:"UserID" bson:"user_id"`
	CreatedAt time.Time          `json:"CreatedAt" bson:"created_at"`
}

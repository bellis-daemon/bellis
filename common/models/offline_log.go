package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OfflineLog struct {
	ID         primitive.ObjectID `bson:"_id"`
	EntityID   primitive.ObjectID `bson:"EntityID"`
	EnvoyTime  time.Time          `bson:"EnvoyTime"`
	EnvoyType  string             `bson:"EnvoyType"`
	OnlineTime time.Time          `bson:"OnlineTime"`
	SentryLogs []SentryLog        `bson:"SentryLogs"`
}

type SentryLog struct {
	SentryName   string    `bson:"SentryName"`
	SentryTime   time.Time `bson:"SentryTime"`
	ErrorMessage string    `bson:"ErrorMessage"`
}

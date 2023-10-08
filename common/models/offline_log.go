package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OfflineLog struct {
	ID         primitive.ObjectID `bson:"_id"`
	EntityID   primitive.ObjectID `bson:"EntityID"`
	EntityName string             `bson:"EntityName"`
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

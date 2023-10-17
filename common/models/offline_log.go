package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OfflineLog struct {
	ID             primitive.ObjectID `bson:"_id"`
	EntityID       primitive.ObjectID `bson:"EntityID"`
	EnvoyType      string             `bson:"EnvoyType"`
	OfflineTime    time.Time          `bson:"OfflineTime"`
	OfflineMessage string             `bson:"OfflineMessage"`
	OnlineTime     time.Time          `bson:"OnlineTime"`
	SentryLogs     []SentryLog        `bson:"SentryLogs"`
}

type SentryLog struct {
	SentryName   string    `bson:"SentryName"`
	SentryTime   time.Time `bson:"SentryTime"`
	ErrorMessage string    `bson:"ErrorMessage"`
}
package models

type EnvoyTelegram struct {
	ID     string `json:"ID" bson:"_id"`
	ChatId int64  `json:"ChatId" bson:"ChatId"`
}

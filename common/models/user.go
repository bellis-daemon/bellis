package models

import (
	"context"
	"time"

	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/models/index"
	"github.com/bellis-daemon/bellis/common/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const SALT = "MONGOUSERSALT"

type User struct {
	ID        primitive.ObjectID `json:"ID" bson:"_id"`
	Email     string             `json:"Email" bson:"Email"`
	Password  string             `json:"Password" bson:"Password"`
	CreatedAt time.Time          `json:"CreatedAt" bson:"CreatedAt"`
	IsVip     bool               `json:"IsVip" bson:"IsVip"`
	Envoy     EnvoyPolicy        `json:"Envoy" bson:"Envoy"`
	Timezone  Timezone           `json:"Timezone" bson:"Timezone"`
}

func NewUser() *User {
	ret := &User{
		ID:        primitive.NewObjectID(),
		Email:     "",
		Password:  "",
		CreatedAt: time.Now(),
		IsVip:     false,
		Envoy: EnvoyPolicy{
			OfflineAlert: false,
			PredictAlert: false,
			Sensitive:    3,
		},
		Timezone: DefaultTimezone,
	}
	return ret
}

func hashPassword(pwd string) string {
	return cryptoo.MD5(SALT + cryptoo.MD5(pwd))
}

func (this *User) CheckPassword(pwd string) bool {
	return hashPassword(pwd) == this.Password
}

func (this *User) SetPassword(ctx context.Context, pwd string) error {
	hpwd := hashPassword(pwd)
	_, err := storage.CUser.UpdateOne(ctx, bson.M{
		"_id": this.ID,
	}, bson.M{
		"$set": bson.M{
			"Password": hpwd,
		},
	})
	if err != nil {
		return err
	}
	this.Password = hpwd
	return nil
}

const (
	IsEnvoyEmail = iota + 1000
	IsEnvoyGotify
	IsEnvoySMS
	IsEnvoyTelegram
	IsEnvoyWebhook
)

type EnvoyPolicy struct {
	PolicyID     primitive.ObjectID `json:"PolicyID" bson:"PolicyID"`
	PolicyType   int                `json:"PolicyType" bson:"PolicyType"`
	OfflineAlert bool               `json:"OfflineAlert" bson:"OfflineAlert"`
	PredictAlert bool               `json:"PredictAlert" bson:"PredictAlert"`
	Sensitive    int                `json:"Sensitive" bson:"Sensitive"`
}

type UserGetter interface {
	User() (*User, error)
}

func init() {
	index.RegistrerIndex(&storage.CUser, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "Email", Value: 1},
			},
		},
	})
}
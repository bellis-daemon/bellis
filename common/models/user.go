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
	// Level user plan level
	// default(free) level: 0
	Level    UserLevel   `json:"Level" bson:"Level"`
	Usage    UserUsage   `json:"Usage" bson:"Usage"`
	Envoy    EnvoyPolicy `json:"Envoy" bson:"Envoy"`
	Timezone Timezone    `json:"Timezone" bson:"Timezone"`
}

type UserLevel uint32

func (this UserLevel) Limit() UserUsage {
	switch this {
	case 0:
		return UserUsage{
			EnvoySMSCount: 10,
			EnvoyCount:    1000,
			EntityCount:   10,
		}
	case 1:
		return UserUsage{
			EnvoySMSCount: 100,
			EnvoyCount:    5000,
			EntityCount:   50,
		}
	default:
		return UserUsage{
			EnvoySMSCount: -1,
			EnvoyCount:    -1,
			EntityCount:   -1,
		}
	}
}

type UserUsage struct {
	EnvoySMSCount int32
	EnvoyCount    int32
	EntityCount   int32
}

func NewUser() *User {
	ret := &User{
		ID:        primitive.NewObjectID(),
		Email:     "",
		Password:  "",
		CreatedAt: time.Now(),
		Level:     0,
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

func (this *User) SetProfile(ctx context.Context, policyType EnvoyPolicyType, policy any) error {
	return storage.MongoUseSession(ctx, func(sessionContext mongo.SessionContext) error {
		// delete old policy
		if this.Envoy.PolicyID != primitive.NilObjectID {
			var coll *mongo.Collection
			switch this.Envoy.PolicyType {
			case IsEnvoyGotify:
				coll = storage.CEnvoyGotify
			case IsEnvoyEmail:
				coll = storage.CEnvoyEmail
			case IsEnvoyWebhook:
				coll = storage.CEnvoyWebhook
			case IsEnvoyTelegram:
				coll = storage.CEnvoyTelegram
			}
			_, err := coll.DeleteOne(sessionContext, bson.M{
				"_id": this.Envoy.PolicyID,
			})
			if err != nil {
				return err
			}
		}
		// create new policy
		switch policyType {
		case IsEnvoyGotify:
			one, err := storage.CEnvoyGotify.InsertOne(sessionContext, policy)
			if err != nil {
				return err
			}
			this.Envoy.PolicyID = one.InsertedID.(primitive.ObjectID)
		case IsEnvoyEmail:
			one, err := storage.CEnvoyEmail.InsertOne(sessionContext, policy)
			if err != nil {
				return err
			}
			this.Envoy.PolicyID = one.InsertedID.(primitive.ObjectID)
		case IsEnvoyWebhook:
			one, err := storage.CEnvoyWebhook.InsertOne(sessionContext, policy)
			if err != nil {
				return err
			}
			this.Envoy.PolicyID = one.InsertedID.(primitive.ObjectID)
		case IsEnvoyTelegram:
			one, err := storage.CEnvoyTelegram.InsertOne(sessionContext, policy)
			if err != nil {
				return err
			}
			this.Envoy.PolicyID = one.InsertedID.(primitive.ObjectID)
		}
		this.Envoy.PolicyType = policyType
		// modify user model
		_, err := storage.CUser.ReplaceOne(sessionContext, bson.M{"_id": this.ID}, this)
		if err != nil {
			return err
		}
		return nil
	})
}

type EnvoyPolicyType int

const (
	IsEnvoyEmail EnvoyPolicyType = iota + 1000
	IsEnvoyGotify
	IsEnvoySMS
	IsEnvoyTelegram
	IsEnvoyWebhook
)

type EnvoyPolicy struct {
	PolicyID     primitive.ObjectID `json:"PolicyID" bson:"PolicyID"`
	PolicyType   EnvoyPolicyType    `json:"PolicyType" bson:"PolicyType"`
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

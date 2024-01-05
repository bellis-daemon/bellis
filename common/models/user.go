package models

import (
	"context"
	"errors"
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
	ID             primitive.ObjectID `json:"ID" bson:"_id"`
	Email          string             `json:"Email" bson:"Email"`
	Password       string             `json:"Password" bson:"Password"`
	CreatedAt      time.Time          `json:"CreatedAt" bson:"CreatedAt"`
	Level          UserLevel          `json:"Level" bson:"Level"` // default(free) level: 0
	LevelExpireAt  time.Time          `json:"LevelExpireAt" bson:"LevelExpireAt"`
	Usage          UserUsage          `json:"Usage" bson:"Usage"`
	Envoy          EnvoyPolicy        `json:"Envoy" bson:"Envoy"`
	Timezone       Timezone           `json:"Timezone" bson:"Timezone"`
	CustomSentries []string           `json:"CustomSentries" bson:"CustomSentries"`
}

type UserLevel int32

const (
	UserLevelAdmin UserLevel = iota - 1
	UserLevelFree
	UserLevelBasic
	UserLevelStandard
	UserLevelPremium
)

func (this UserLevel) Limit() UserUsage {
	switch this {
	case UserLevelAdmin:
		return UserUsage{
			EnvoySMSCount:    -1,
			EnvoyCount:       -1,
			EntityCount:      -1,
			EnvoyPolicyCount: -1,
		}
	case UserLevelFree:
		return UserUsage{
			EnvoySMSCount:    5,
			EnvoyCount:       100,
			EntityCount:      5,
			EnvoyPolicyCount: 1,
		}
	case UserLevelBasic:
		return UserUsage{
			EnvoySMSCount:    50,
			EnvoyCount:       5000,
			EntityCount:      30,
			EnvoyPolicyCount: 5,
		}
	case UserLevelStandard:
		return UserUsage{
			EnvoySMSCount:    200,
			EnvoyCount:       10000,
			EntityCount:      60,
			EnvoyPolicyCount: 10,
		}
	case UserLevelPremium:
		return UserUsage{
			EnvoySMSCount:    1000,
			EnvoyCount:       50000,
			EntityCount:      200,
			EnvoyPolicyCount: 30,
		}
	}
	panic("invalid user level")
}

type UserUsage struct {
	EnvoySMSCount    int32 `json:"EnvoySMSCount" bson:"EnvoySMSCount"`
	EnvoyCount       int32 `json:"EnvoyCount" bson:"EnvoyCount"`
	EntityCount      int32 `json:"EntityCount" bson:"EntityCount"`
	EnvoyPolicyCount int32 `json:"EnvoyPolicyCount" bson:"EnvoyPolicyCount"`
}

func NewUser() *User {
	ret := &User{
		ID:            primitive.NewObjectID(),
		Email:         "",
		Password:      "",
		CreatedAt:     time.Now(),
		Level:         UserLevelFree,
		LevelExpireAt: time.Time{},
		Envoy: EnvoyPolicy{
			OfflineAlert: false,
			Sensitive:    3,
		},
		Usage: UserUsage{
			EnvoySMSCount:    0,
			EnvoyCount:       0,
			EntityCount:      0,
			EnvoyPolicyCount: 0,
		},
		Timezone: DefaultTimezone,
	}
	return ret
}

func hashPassword(pwd string) string {
	return cryptoo.MD5(SALT + cryptoo.MD5(pwd))
}

func (this *User) SetUserLevel(ctx context.Context, level UserLevel, ttl ...time.Duration) error {
	this.Level = level
	this.LevelExpireAt = time.Time{}
	if len(ttl) > 0 {
		this.LevelExpireAt = time.Now().Add(ttl[0])
	}
	_, err := storage.CUser.UpdateByID(ctx, this.ID, bson.M{"$set": bson.M{"Level": this.Level, "LevelExpireAt": this.LevelExpireAt}})
	return err
}

func (this *User) UsageEnvoySMSAccessible() bool {
	return this.Usage.EnvoySMSCount < this.Level.Limit().EnvoySMSCount
}

func (this *User) UsageEnvoyAccessible() bool {
	return this.Usage.EnvoyCount < this.Level.Limit().EnvoyCount
}

func (this *User) UsageEntityAccessible() bool {
	return this.Usage.EntityCount < this.Level.Limit().EntityCount
}

func (this *User) UsageEnvoyPolicyAccessible() bool {
	return this.Usage.EnvoyPolicyCount < this.Level.Limit().EnvoyPolicyCount
}

func (this *User) UsageEnvoySMSIncr(ctx context.Context, delta int32) error {
	this.Usage.EnvoySMSCount += delta
	_, err := storage.CUser.UpdateByID(ctx, this.ID, bson.M{"$set": bson.M{"Usage.EnvoySMSCount": this.Usage.EnvoySMSCount}})
	return err
}

func (this *User) UsageEnvoyIncr(ctx context.Context, delta int32) error {
	this.Usage.EnvoyCount += delta
	_, err := storage.CUser.UpdateByID(ctx, this.ID, bson.M{"$set": bson.M{"Usage.EnvoyCount": this.Usage.EnvoyCount}})
	return err
}

func (this *User) UsageEntityIncr(ctx context.Context, delta int32) error {
	this.Usage.EntityCount += delta
	_, err := storage.CUser.UpdateByID(ctx, this.ID, bson.M{"$set": bson.M{"Usage.EntityCount": this.Usage.EntityCount}})
	return err
}

func (this *User) UsageEnvoyPolicyIncr(ctx context.Context, delta int32) error {
	this.Usage.EnvoyPolicyCount += delta
	_, err := storage.CUser.UpdateByID(ctx, this.ID, bson.M{"$set": bson.M{"Usage.EnvoyPolicyCount": this.Usage.EnvoyPolicyCount}})
	return err
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
			default:
				return errors.New("invalid policy type")
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
		default:
			return errors.New("invalid policy type")
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

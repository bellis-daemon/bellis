package models

import (
	"context"
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const SALT = "MONGOUSERSALT"

type User struct {
	ID        primitive.ObjectID `json:"ID" bson:"_id"`
	Email     string             `json:"Email" bson:"Email"`
	Password  string             `json:"Password" bson:"Password"`
	CreatedAt time.Time          `json:"CreatedAt" bson:"CreatedAt"`
	IsVip     bool               `json:"IsVip" bson:"IsVip"`
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

package consumer

import (
	"context"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func registerEmailCaptcha() {
	instance().Register("CaptchaToEmail", func(message *redistream.Message) error {
		glgf.Debug(message)
		return nil
	})
}

func registerEntityAlert() {
	instance().Register("EntityAlert", func(message *redistream.Message) error {
		id, err := primitive.ObjectIDFromHex(message.Values["EntityID"].(string))
		if err != nil {
			return err
		}
		var entity models.Application
		err = storage.CEntity.FindOne(context.Background(), bson.M{"_id": id}).Decode(&entity)

		return nil
	})
}

func init() {
	glgf.Debug("register EmailCaptcha")
	registerEmailCaptcha()
	registerEntityAlert()
}

package consumer

import (
	"context"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/email"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/gotify"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func entityOfflineAlert() {
	redistream.Instance().Register("EntityOfflineAlert", func(message *redistream.Message) error {
		glgf.Debug(message)
		id, err := primitive.ObjectIDFromHex(message.Values["EntityID"].(string))
		if err != nil {
			return err
		}
		var entity models.Application
		err = storage.CEntity.FindOne(context.Background(), bson.M{"_id": id}).Decode(&entity)
		if err != nil {
			return err
		}
		var user models.User
		err = storage.CUser.FindOne(context.Background(), bson.M{"_id": entity.UserID}).Decode(&user)
		if err != nil {
			return err
		}
		switch user.Envoy.PolicyType {
		case models.IsEnvoyGotify:
			var policy models.EnvoyGotify
			err = storage.CEnvoyGotify.FindOne(context.Background(), bson.M{"_id": user.Envoy.PolicyID}).Decode(&policy)
			if err != nil {
				return err
			}
			err = gotify.AlertOffline(&entity, &policy, message.Values["Message"].(string))
			if err != nil {
				return err
			}
		case models.IsEnvoyEmail:
			var policy models.EnvoyEmail
			err = storage.CEnvoyEmail.FindOne(context.Background(), bson.M{"_id": user.Envoy.PolicyID}).Decode(&policy)
			if err != nil {
				return err
			}
			err = email.AlertOffline(&entity, &policy, message.Values["Message"].(string))
			if err != nil {
				return err
			}
		default:
			glgf.Warn("User envoy policy is empty, ignoring", entity.Name, user.Envoy)
			return nil
		}
		return nil
	})
}

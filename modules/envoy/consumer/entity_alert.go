package consumer

import (
	"context"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/email"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/gotify"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/webhook"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func entityOfflineAlert() {
	stream.Register(common.EntityOfflineAlert, func(ctx context.Context, message *redistream.Message) error {
		glgf.Debug(message)
		offlineTime := time.UnixMilli(cast.ToInt64(message.Values["OfflineTime"]))
		id, err := primitive.ObjectIDFromHex(message.Values["EntityID"].(string))
		if err != nil {
			return err
		}
		var entity models.Application
		err = storage.CEntity.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
		if err != nil {
			return err
		}
		var user models.User
		err = storage.CUser.FindOne(ctx, bson.M{"_id": entity.UserID}).Decode(&user)
		if err != nil {
			return err
		}
		envoyType := ""
		switch user.Envoy.PolicyType {
		case models.IsEnvoyGotify:
			envoyType = "Gotify"
			err = gotify.New(ctx).WithPolicyId(user.Envoy.PolicyID).AlertOffline(&entity, message.Values["Message"].(string), offlineTime)
			if err != nil {
				return err
			}
		case models.IsEnvoyEmail:
			envoyType = "Email"
			err = email.New(ctx).WithPolicyId(user.Envoy.PolicyID).AlertOffline(&entity, message.Values["Message"].(string), offlineTime)
			if err != nil {
				return err
			}
		case models.IsEnvoyWebhook:
			envoyType = "Webhook"
			err = webhook.New(ctx).WithPolicyId(user.Envoy.PolicyID).AlertOffline(&entity, message.Values["Message"].(string), offlineTime)
		default:
			glgf.Warn("User envoy policy is empty, ignoring", entity.Name, user.Envoy)
			return nil
		}
		go func() {
			retry.Do(func() error {
				err := writeOfflineLog(ctx, &entity, offlineTime, envoyType)
				return err
			}, retry.Context(ctx))
		}()
		return nil
	})
}

func writeOfflineLog(ctx context.Context, entity *models.Application, offlineTime time.Time, envoyType string) error {
	log := models.OfflineLog{
		ID:         primitive.NewObjectID(),
		EntityID:   entity.ID,
		EnvoyTime:  time.Now(),
		EnvoyType:  envoyType,
		SentryLogs: []models.SentryLog{},
	}
	query, err := storage.QueryInfluxDB.Query(ctx, fmt.Sprintf(`
from(bucket: "backend")
  |> range(start: -5m, stop: %s)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")
  |> filter(fn: (r) => r["_field"] == "c_err" or r["_field"] == "c_sentry")
  |> sort(columns: ["_time"], desc: true)
  |> limit(n: 3)
  |> sort(columns: ["_time"], desc: false)
  |> group(columns: ["_time"])
`, offlineTime.Format(time.RFC3339), common.Measurements[entity.SchemeID], entity.ID.Hex()))
	if err != nil {
		return err
	}
	for query.Next() {
		sl := models.SentryLog{
			SentryTime: query.Record().Time(),
		}
		if query.Record().Field() == "c_err" {
			sl.ErrorMessage = cast.ToString(query.Record().Value())
		}
		if query.Record().Field() == "c_sentry" {
			sl.SentryName = cast.ToString(query.Record().Value())
		}
		query.Next()
		if query.Record().Field() == "c_err" {
			sl.ErrorMessage = cast.ToString(query.Record().Value())
		}
		if query.Record().Field() == "c_sentry" {
			sl.SentryName = cast.ToString(query.Record().Value())
		}
		log.SentryLogs = append(log.SentryLogs, sl)
	}
	glgf.Debug(log)
	_, err = storage.COfflineLog.InsertOne(ctx, log)
	return err
}

func entityOnlineAlert() {
	stream.Register(common.EntityOnlineAlert, func(ctx context.Context, message *redistream.Message) error {
		glgf.Debug(message)
		onlineTime := time.UnixMilli(cast.ToInt64(message.Values["OfflineTime"]))
		id, err := primitive.ObjectIDFromHex(message.Values["EntityID"].(string))
		if err != nil {
			return err
		}
		var entity models.Application
		err = storage.CEntity.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
		if err != nil {
			return err
		}
		err = retry.Do(func() error {
			return writeOnlineLog(ctx, &entity, onlineTime)
		})
		if err != nil {
			glgf.Error(err)
		}
		return err
	})
}

func writeOnlineLog(ctx context.Context, entity *models.Application, onlineTIme time.Time) error {
	var log models.OfflineLog
	err := storage.COfflineLog.FindOne(ctx, bson.M{"EntityID": entity.ID}, options.FindOne().SetSort(bson.M{"$natural": -1})).Decode(&log)
	if err != nil {
		return err
	}
	_, err = storage.COfflineLog.UpdateOne(ctx, bson.M{"_id": log.ID}, bson.M{"$set": bson.M{"OnlineTime": onlineTIme}})
	return err
}

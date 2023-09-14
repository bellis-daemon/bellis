package consumer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/email"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/gotify"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/telegram"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/webhook"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func entityOfflineAlert() {
	stream.Register(common.EntityOfflineAlert, func(ctx context.Context, message *redistream.Message) error {
		offlineTime := time.UnixMilli(cast.ToInt64(message.Values["OfflineTime"]))
		id, err := primitive.ObjectIDFromHex(message.Values["EntityID"].(string))
		if err != nil {
			return fmt.Errorf("cant parse hex user id: %s: %w", message.Values["EntityID"].(string), err)
		}
		var entity models.Application
		err = storage.CEntity.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				glgf.Errorf("cant find entity using entity id: %s: %s", id.Hex(), err.Error())
				return nil
			}
			return fmt.Errorf("finding entity mongo inernal error: %w", err)
		}
		var user models.User
		err = storage.CUser.FindOne(ctx, bson.M{"_id": entity.UserID}).Decode(&user)
		if err != nil {
			return fmt.Errorf("cant find user using user id: %s: %w", entity.UserID.Hex(), err)
		}
		// check if entity is previously offline
		ok, err := isOnlineState(ctx, &entity)
		if err != nil {
			return fmt.Errorf("cant get entity offline log: %w", err)
		}
		if !ok {
			glgf.Warn("entity alert canceled because of previously offline: ", entity, message)
			return nil
		}
		envoyType := ""
		var envoyDriver drivers.EnvoyDriver
		switch user.Envoy.PolicyType {
		case models.IsEnvoyGotify:
			envoyType = "Gotify"
			envoyDriver = gotify.New(ctx).WithPolicyId(user.Envoy.PolicyID)
		case models.IsEnvoyEmail:
			envoyType = "Email"
			envoyDriver = email.New(ctx).WithPolicyId(user.Envoy.PolicyID)
		case models.IsEnvoyWebhook:
			envoyType = "Webhook"
			envoyDriver = webhook.New(ctx).WithPolicyId(user.Envoy.PolicyID)
		case models.IsEnvoyTelegram:
			envoyType = "Telegram"
			envoyDriver = telegram.New(ctx).WithPolicyId(user.Envoy.PolicyID)
		default:
			glgf.Warn("User envoy policy is empty, ignoring", entity.Name, user.Envoy)
			return nil
		}
		err = envoyDriver.AlertOffline(&entity, message.Values["Message"].(string), offlineTime)
		if err != nil {
			return fmt.Errorf("cant find policy using policy id: %s, %w", user.Envoy.PolicyID.Hex(), err)
		}
		err = retry.Do(func() error {
			err := writeOfflineLog(ctx, &entity, offlineTime, envoyType)
			return err
		}, retry.Context(ctx), retry.Attempts(3))
		if err != nil {
			return err
		}
		glgf.Debug("Offline alert sent: ", entity.Name)
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
  |> range(start: %s, stop: %s)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")
  |> filter(fn: (r) => r["_field"] == "c_err" or r["_field"] == "c_sentry")
  |> sort(columns: ["_time"], desc: true)
  |> limit(n: 3)
  |> sort(columns: ["_time"], desc: false)
  |> group(columns: ["_time"])
`, offlineTime.Add(-5*time.Minute).Format(time.RFC3339), offlineTime.Format(time.RFC3339), entity.Scheme, entity.ID.Hex()))
	if err != nil {
		return fmt.Errorf("error querying influxdb: %w", err)
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
	if err != nil {
		return fmt.Errorf("error inserting offline log: %w", err)
	}
	return nil
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
			if errors.Is(err, mongo.ErrNoDocuments) {
				glgf.Errorf("cant find entity using entity id: %s: %s", id.Hex(), err.Error())
				return nil
			}
			return fmt.Errorf("finding entity mongo inernal error: %w", err)
		}
		err = retry.Do(func() error {
			return writeOnlineLog(ctx, &entity, onlineTime)
		}, retry.Context(ctx), retry.Attempts(3))
		if err != nil {
			return err
		}
		return nil
	})
}

func writeOnlineLog(ctx context.Context, entity *models.Application, onlineTIme time.Time) error {
	var log models.OfflineLog
	err := storage.COfflineLog.FindOne(ctx, bson.M{"EntityID": entity.ID}, options.FindOne().SetSort(bson.M{"$natural": -1})).Decode(&log)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			glgf.Errorf("cant find offline log using entity id: %s: %s", entity.ID.Hex(), err.Error())
			return nil
		}
		return fmt.Errorf("finding entity offline log inernal error: %w", err)
	}
	_, err = storage.COfflineLog.UpdateOne(ctx, bson.M{"_id": log.ID}, bson.M{"$set": bson.M{"OnlineTime": onlineTIme}})
	if err != nil {
		return fmt.Errorf("error updating offline log in mongodb: %w", err)
	}
	return nil
}

func isOnlineState(ctx context.Context, entity *models.Application) (bool, error) {
	var log models.OfflineLog
	err := storage.COfflineLog.FindOne(ctx, bson.M{"EntityID": entity.ID}, options.FindOne().SetSort(bson.M{"$natural": -1})).Decode(&log)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return false, fmt.Errorf("internal mongodb err: %w", err)
		} else {
			return true, nil
		}
	}
	return !log.OnlineTime.IsZero(), nil
}

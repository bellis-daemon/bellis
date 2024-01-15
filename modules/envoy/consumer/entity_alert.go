package consumer

import (
	"context"
	"errors"
	"fmt"
	"sync"
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
		user, err := entity.User()
		if err != nil {
			return fmt.Errorf("cant find user using user id: %s: %w", entity.UserID.Hex(), err)
		}

		// check if entity is previously offline
		var log *models.OfflineLog
		log, err = recentOfflineLog(ctx, &entity)
		if err != nil {
			return fmt.Errorf("cant get entity offline log: %w", err)
		}
		if !log.OnlineTime.IsZero() {
			log, err = writeOfflineLog(ctx, &entity, message.Values["Message"].(string), offlineTime)
			if err != nil {
				return fmt.Errorf("cant write offline log: %w", err)
			}
		} else {
			glgf.Warn("Offline alert message received for an already offline entity, ignoring: %s", entity.ID.Hex())
			return nil
		}
		var wg sync.WaitGroup
		failed := false
		for i := range user.EnvoyPolicies {
			policyId := user.EnvoyPolicies[i].PolicyID
			policyType := user.EnvoyPolicies[i].PolicyType
			wg.Add(1)
			go func() {
				defer wg.Done()
				if policyType == models.IsEnvoySMS && !user.UsageEnvoySMSAccessible() {
					glgf.Warn("User <%s>`s envoy sms usage exceeds: %d.", user.Usage.EnvoySMSCount)
					return
				}
				if policyType != models.IsEnvoySMS && !user.UsageEnvoyAccessible() {
					glgf.Warn("User <%s>`s envoy usage exceeds: %d.", user.Usage.EnvoyCount)
					return
				}
				count, err := storage.CEnvoyLog.CountDocuments(ctx, bson.M{
					"$and": bson.D{
						{
							"OfflineLogID", log.ID,
						},
						{
							"Success", true,
						},
						{
							"PolicySnapShot.ID", policyId,
						},
					},
				})
				if err != nil {
					failed = true
					glgf.Warn(err)
					return
				}
				if count != 0 {
					glgf.Info("Previous alert message success via current policy, ignoring.", policyId.String(), policyType)
					return
				}
				envoyType := ""
				var envoyDriver drivers.EnvoyDriver
				switch policyType {
				case models.IsEnvoyGotify:
					envoyType = "Gotify"
					envoyDriver = gotify.New(ctx).WithPolicyId(policyId)
				case models.IsEnvoyEmail:
					envoyType = "Email"
					envoyDriver = email.New(ctx).WithPolicyId(policyId)
				case models.IsEnvoyWebhook:
					envoyType = "Webhook"
					envoyDriver = webhook.New(ctx).WithPolicyId(policyId)
				case models.IsEnvoyTelegram:
					envoyType = "Telegram"
					envoyDriver = telegram.New(ctx).WithPolicyId(policyId)
				default:
					glgf.Warn("User envoy policy is empty, ignoring", entity.Name, policyId.String(), policyType)
					failed = true
					return
				}
				err = envoyDriver.AlertOffline(user, &entity, log)
				{
					envoyLog := &models.EnvoyLog{
						ID:             primitive.NewObjectID(),
						SendTime:       time.Now(),
						Success:        err == nil,
						OfflineLogID:   log.ID,
						PolicyType:     envoyType,
						PolicySnapShot: envoyDriver.PolicySnapShot(),
					}
					if err != nil {
						envoyLog.FailedMessage = err.Error()
					}
					_, err := storage.CEnvoyLog.InsertOne(context.Background(), envoyLog)
					if err != nil {
						glgf.Error(err)
					}
				}
				if err != nil {
					glgf.Error("send offline alert failed: %w", err)
					failed = true
					return
				}
				go func() {
					if policyType == models.IsEnvoySMS {
						err = user.UsageEnvoySMSIncr(ctx, 1)
						if err != nil {
							glgf.Error(err)
						}
					} else {
						err = user.UsageEnvoyIncr(ctx, 1)
						if err != nil {
							glgf.Error(err)
						}
					}
				}()
				glgf.Debugf("Offline alert of %s sent via %s", entity.Name, envoyType)
			}()
		}
		wg.Wait()
		if failed == true {
			return errors.New("error while sending alert")
		}
		return nil
	})
}

func writeOfflineLog(ctx context.Context, entity *models.Application, envoyMessage string, offlineTime time.Time) (*models.OfflineLog, error) {
	log := &models.OfflineLog{
		ID:             primitive.NewObjectID(),
		EntityID:       entity.ID,
		OfflineTime:    time.Now(),
		OfflineMessage: envoyMessage,
		SentryLogs:     []models.SentryLog{},
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
		return nil, fmt.Errorf("error querying influxdb: %w", err)
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
		return nil, fmt.Errorf("error inserting offline log: %w", err)
	}
	return log, nil
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
			return nil
		}
		return fmt.Errorf("finding entity offline log inernal error: %w", err)
	}
	if !log.OnlineTime.IsZero() {
		return nil
	}
	_, err = storage.COfflineLog.UpdateOne(ctx, bson.M{"_id": log.ID}, bson.M{"$set": bson.M{"OnlineTime": onlineTIme}})
	if err != nil {
		return fmt.Errorf("error updating offline log in mongodb: %w", err)
	}
	return nil
}

func recentOfflineLog(ctx context.Context, entity *models.Application) (*models.OfflineLog, error) {
	var log models.OfflineLog
	err := storage.COfflineLog.FindOne(ctx, bson.M{"EntityID": entity.ID}, options.FindOne().SetSort(bson.M{"$natural": -1})).Decode(&log)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("internal mongodb err: %w", err)
		} else {
			return &models.OfflineLog{OnlineTime: time.Now()}, nil
		}
	}
	return &log, nil
}

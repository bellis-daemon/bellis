package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bellis-daemon/bellis/common/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
)

func setTelegramWebhook() {
	if storage.Config().TelegramBotToken != "" {
		webhookEndpoint := storage.Config().WebEndpoint + "/api/callback/telegram"
		target := storage.Config().TelegramBotApiEndpoint + fmt.Sprintf("/bot%s/setWebhook", storage.Config().TelegramBotToken)
		body := map[string]any{
			"url": webhookEndpoint,
			"allowed_updates": []string{
				"message",
			},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			glgf.Error("err while setting telegram webhook", err)
			return
		}
		req, err := http.NewRequest("POST", target, &buf)
		if err != nil {
			glgf.Error("err while setting telegram webhook", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		glgf.Debug(webhookEndpoint)
		get, err := http.DefaultClient.Do(req)
		if err != nil {
			glgf.Error("err while setting telegram webhook", err)
			return
		}
		buf.Reset()
		buf.ReadFrom(get.Body)
		glgf.Info("updated telegram webhook", get.Status, buf.String())
	}
}

func resetUserUsages() {
	ctx := context.Background()
	many, err := storage.CUser.UpdateMany(ctx, bson.M{}, bson.M{
		"$set": bson.M{
			"Usage.EnvoySMSCount":    0,
			"Usage.EnvoyCount":       0,
			"Usage.EnvoyPolicyCount": 0,
		},
	})
	if err != nil {
		glgf.Error(err)
		return
	}
	glgf.Infof("Reseted user usages, <%d> users matched, <%d> users modified.", many.MatchedCount, many.ModifiedCount)
}

func checkUserEntityUsageCount() {
	ctx := context.Background()
	find, err := storage.CUser.Find(ctx, bson.M{})
	if err != nil {
		glgf.Error(err)
		return
	}
	var user models.User
	for find.Next(ctx) {
		err := find.Decode(&user)
		if err != nil {
			glgf.Error(err)
			continue
		}
		count, err := storage.CEntity.CountDocuments(ctx, bson.M{"UserID": user.ID})
		if err != nil {
			glgf.Error(err)
			continue
		}
		if int32(count) == user.Usage.EntityCount {
			continue
		}
		updated, err := storage.CUser.UpdateByID(ctx, user.ID, bson.M{
			"$set": bson.M{
				"Usage.EntityCount": count,
			},
		})
		if err != nil {
			glgf.Error(err)
			return
		}
		glgf.Infof("User <%s> has wrong entity usage count, updated: %d => %d,modified document: %d", user.Email, user.Usage.EntityCount, count, updated.ModifiedCount)
	}
}

func checkUserLevelExpire() {
	ctx := context.Background()
	find, err := storage.CUser.Find(ctx, bson.M{"$and": bson.A{
		bson.D{{
			Key:   "LevelExpireAt",
			Value: bson.E{Key: "$lt", Value: time.Now()},
		}},
		bson.D{{
			Key:   "LevelExpireAt",
			Value: bson.E{Key: "$ne", Value: time.Time{}},
		}},
	}})
	if err != nil {
		glgf.Error(err)
		return
	}
	var user models.User
	for find.Next(ctx) {
		err := find.Decode(&user)
		if err != nil {
			glgf.Error(err)
			continue
		}
		glgf.Infof("User <%s>`s level expires (%s), moving user to free level.", user.ID, user.LevelExpireAt.Format(time.DateTime))
		user.LevelExpireAt = time.Time{}
		user.Level = models.UserLevelFree
		_, err = storage.CUser.ReplaceOne(ctx, bson.M{"_id": user.ID}, &user)
		if err != nil {
			glgf.Error(err)
			continue
		}
		if user.Usage.EntityCount > user.Level.Limit().EntityCount {
			//todo: do something
		}
	}
}

package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bellis-daemon/bellis/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"

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
	glgf.Debug("reseting user usages")
	//todo: implement reset usage function
}

func checkUserLevelExpire() {
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
		if user.LevelExpireAt.IsZero() {
			continue
		}
		if user.LevelExpireAt.After(time.Now()) {
			continue
		}
		user.LevelExpireAt = time.Time{}
		user.Level = models.UserLevelFree
		_, err = storage.CUser.ReplaceOne(ctx, bson.M{"_id": user.ID}, &user)
		if err != nil {
			glgf.Error(err)
			continue
		}
	}
}

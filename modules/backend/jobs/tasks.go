package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
	"net/http"
	"time"
)

func clearInfluxdbDataExpired() {
	glgf.Info("deleting data before one week in influxdb")
	err := storage.DeleteInfluxDB.DeleteWithName(context.Background(), "bellis", "backend", time.UnixMilli(0), time.Now().AddDate(0, 0, -7), "")
	if err != nil {
		glgf.Error(err)
	}
}

func setTelegramWebhook() {
	if storage.Config().TelegramBotToken != "" {
		webhookEndpoint := storage.Config().WebEndpoint + "/callback/telegram"
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

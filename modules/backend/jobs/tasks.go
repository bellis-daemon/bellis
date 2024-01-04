package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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
